#!/usr/bin/env python3
# pylint: disable=invalid-name
# Description: Builds a standalone copy of binutils
# 
# Original: https://github.com/ClangBuiltLinux/tc-build
# License: Apachev2
#
# Edited by: Luciano Remes
#
# Updated build settings and CFLAGS 
######

import argparse
import hashlib
import multiprocessing
import pathlib
import platform
import re
import sys
import shutil
import subprocess
import time

BINUTIL_VERSION = "binutils-2.40"

# UTILS:
def create_gitignore(folder):
    """
    Create a gitignore that ignores all files in a folder. Some folders are not
    known until the script is run so they can't be added to the root .gitignore
    :param folder: Folder to create the gitignore in
    """
    folder.joinpath('gitignore').write_text('*\n', encoding='utf-8')

def libc_is_musl():
    """
    Returns whether or not the current libc is musl or not.
    """
    # musl's ldd does not appear to support '--version' directly, as its return
    # code is 1 and it prints all text to stderr. However, it does print the
    # version information so it is good enough. Just 'check=False' it and move
    # on.
    ldd_out = subprocess.run(['ldd', '--version'],
                             capture_output=True,
                             check=False,
                             text=True)
    if re.search('musl', ldd_out.stderr if ldd_out.stderr else ldd_out.stdout):
        return True
    return False

def get_duration(start_seconds, end_seconds=None):
    """
    Formats a duration in days, hours, minutes, and seconds.
    :param start_seconds: The start of the duration
    :param end_seconds: The end of the duration; can be omitted for current time
    :return: A string with the non-zero parts of the duration.
    """
    if not end_seconds:
        end_seconds = time.time()
    seconds = int(end_seconds - start_seconds)
    days, seconds = divmod(seconds, 60 * 60 * 24)
    hours, seconds = divmod(seconds, 60 * 60)
    minutes, seconds = divmod(seconds, 60)

    parts = []
    if days:
        parts.append(f"{days}d")
    if hours:
        parts.append(f"{hours}h")
    if minutes:
        parts.append(f"{minutes}m")
    parts.append(f"{seconds}s")

    return ' '.join(parts)

def flush_std_err_out():
    sys.stderr.flush()
    sys.stdout.flush()

def print_header(string):
    """
    Prints a fancy header
    :param string: String to print inside the header
    """
    # Use bold cyan for the header so that the headers
    # are not intepreted as success (green) or failed (red)
    print("\033[01;36m")
    for _ in range(0, len(string) + 6):
        print("=", end="")
    print(f"\n== {string} ==")
    for _ in range(0, len(string) + 6):
        print("=", end="")
    # \033[0m resets the color back to the user's default
    print("\n\033[0m")
    flush_std_err_out()

def download_binutils(folder):
    """
    Downloads the latest stable version of binutils
    :param folder: Directory to download binutils to
    """
    binutils = BINUTIL_VERSION
    binutils_folder = folder.joinpath(binutils)
    if not binutils_folder.is_dir():
        # Remove any previous copies of binutils
        for entity in folder.glob('binutils-*'):
            if entity.is_dir():
                shutil.rmtree(entity)
            else:
                entity.unlink()

        # Download the tarball
        binutils_tarball = folder.joinpath(binutils + ".tar.xz")
        curl_cmd = [
            "curl", "-LSs", "-o", binutils_tarball,
            f"https://sourceware.org/pub/binutils/releases/{binutils_tarball.name}"
        ]
        subprocess.run(curl_cmd, check=True)
        verify_binutils_checksum(binutils_tarball)
        # Extract the tarball then remove it
        subprocess.run(["tar", "-xJf", binutils_tarball.name],
                       check=True,
                       cwd=folder)
        create_gitignore(binutils_folder)
        binutils_tarball.unlink()


def verify_binutils_checksum(file_to_check):
    # Check the SHA512 checksum of the downloaded file with a known good one
    file_hash = hashlib.sha512()
    with file_to_check.open("rb") as file:
        while True:
            data = file.read(131072)
            if not data:
                break
            file_hash.update(data)
    # Get good hash from file
    curl_cmd = [
        'curl', '-fLSs',
        'https://sourceware.org/pub/binutils/releases/sha512.sum'
    ]
    sha512_sums = subprocess.run(curl_cmd,
                                 capture_output=True,
                                 check=True,
                                 text=True).stdout
    line_match = fr"([0-9a-f]+)\s+{file_to_check.name}$"
    if not (match := re.search(line_match, sha512_sums, flags=re.M)):
        raise RuntimeError(
            "Could not find binutils hash in sha512.sum output?")
    if file_hash.hexdigest() != match.groups()[0]:
        raise RuntimeError(
            "binutils: SHA512 checksum does not match known good one!")


def host_arch_target():
    """
    Converts the host architecture to the first part of a target triple
    :return: Target host
    """
    host_mapping = {
        "armv7l": "arm",
        "ppc64": "powerpc64",
        "ppc64le": "powerpc64le",
        "ppc": "powerpc"
    }
    machine = platform.machine()
    return host_mapping.get(machine, machine)


def target_arch(target):
    """
    Returns the architecture from a target triple
    :param target: Triple to deduce architecture from
    :return: Architecture associated with given triple
    """
    return target.split("-")[0]


def host_is_target(target):
    """
    Checks if the current target triple the same as the host.
    :param target: Triple to match host architecture against
    :return: True if host and target are same, False otherwise
    """
    return host_arch_target() == target_arch(target)


def parse_parameters(root_folder):
    """
    Parses parameters passed to the script into options
    :param root_folder: The directory where the script is being invoked from
    :return: A 'Namespace' object with all the options parsed from supplied parameters
    """
    parser = argparse.ArgumentParser()
    parser.add_argument("-b",
                        "--binutils-folder",
                        help="""
                        By default, the script will download a copy of the binutils source in the same folder as
                        this script. If you have your own copy of the binutils source that you would like to build
                        from, pass it to this parameter. This can either be an absolute or relative path.
                        """,
                        type=str)
    parser.add_argument("-B",
                        "--build-folder",
                        help="""
                        By default, the script will create a "build" folder in the same folder as this script,
                        then a "binutils" folder within that one and build the files there. If you would like
                        that done somewhere else, pass it to this parameter. This can either be an absolute
                        or relative path.
                        """,
                        type=str,
                        default=root_folder.joinpath("build", "binutils"))
    parser.add_argument("-I",
                        "--install-folder",
                        help="""
                        By default, the script will create an "install" folder in the same folder as this script
                        and install binutils there. If you'd like to have it installed somewhere else, pass
                        it to this parameter. This can either be an absolute or relative path.
                        """,
                        type=str,
                        default=root_folder.joinpath("install"))
    parser.add_argument("-s",
                        "--skip-install",
                        help="""
                        Skip installing binutils into INSTALL_FOLDER
                        """,
                        action="store_true")
    parser.add_argument("-t",
                        "--targets",
                        help="""
                        The script can build binutils targeting arm-linux-gnueabi, aarch64-linux-gnu,
                        mipsel-linux-gnu, powerpc-linux-gnu, powerpc64-linux-gnu, powerpc64le-linux-gnu,
                        riscv64-linux-gnu, s390x-linux-gnu, and x86_64-linux-gnu.

                        You can either pass the full target or just the first part (arm, aarch64, x86_64, etc)
                        or all if you want to build all targets (which is the default). It will only add the
                        target prefix if it is not for the host architecture.
                        """,
                        nargs="+")
    parser.add_argument("-m",
                        "--march",
                        metavar="ARCH",
                        help="""
                        Add -march=ARCH and -mtune=ARCH to CFLAGS to optimize the toolchain for the target
                        host processor.
                        """,
                        type=str)
    return parser.parse_args()


def create_targets(targets):
    """
    Generate a list of targets that can be passed to the binutils compile function
    :param targets: A list of targets to convert to binutils target triples
    :return: A list of target triples
    """
    targets_dict = {
        "arm": "arm-linux-gnueabi",
        "aarch64": "aarch64-linux-gnu",
        "mips": "mips-linux-gnu",
        "mipsel": "mipsel-linux-gnu",
        "powerpc64": "powerpc64-linux-gnu",
        "powerpc64le": "powerpc64le-linux-gnu",
        "powerpc": "powerpc-linux-gnu",
        "riscv64": "riscv64-linux-gnu",
        "s390x": "s390x-linux-gnu",
        "x86_64": "x86_64-linux-gnu"
    }

    targets_set = set()
    for target in targets:
        if target == "all":
            return list(targets_dict.values())
        if target == "host":
            key = host_arch_target()
        else:
            key = target_arch(target)
        targets_set.add(targets_dict[key])

    return list(targets_set)


def cleanup(build_folder):
    """
    Cleanup the build directory
    :param build_folder: Build directory
    """
    if build_folder.is_dir():
        shutil.rmtree(build_folder)
    build_folder.mkdir(parents=True, exist_ok=True)


def invoke_configure(binutils_folder, build_folder, install_folder, target,
                     host_arch):
    """
    Invokes the configure script to generate a Makefile
    :param binutils_folder: Binutils source folder
    :param build_folder: Build directory
    :param install_folder: Directory to install binutils to
    :param target: Target to compile for
    :param host_arch: Host architecture to optimize for
    """
    #CHANGE HERE:
    # CC=clang-15
    configure = [
        binutils_folder.joinpath("configure"),
        'CC=clang',
        'CXX=g++',
        '--disable-compressed-debug-sections',
        '--disable-gdb',
        '--disable-nls',
        '--disable-werror',
        '--enable-deterministic-archives',
        '--enable-new-dtags',
        '--enable-plugins',
        '--enable-threads',
        '--quiet',
        '--with-system-zlib',
    ]  # yapf: disable
    if install_folder:
        configure += [f'--prefix={install_folder}']
    if host_arch:
        configure += [
            f'CFLAGS=-g -O0 -march={host_arch} -mtune={host_arch} -fsanitize=address -fno-omit-frame-pointer -fsanitize-coverage=edge,trace-pc-guard',
            f'CXXFLAGS=-O2 -march={host_arch} -mtune={host_arch}'
        ]
    else:
        configure += ['CFLAGS=-g -O0 -fsanitize=address -fno-omit-frame-pointer -fsanitize-coverage=edge,trace-pc-guard', 'CXXFLAGS=-O2']
    # gprofng uses glibc APIs that might not be available on musl
    if libc_is_musl():
        configure += ['--disable-gprofng']

    configure_arch_flags = {
        "arm-linux-gnueabi": [
            '--disable-multilib',
            '--with-gnu-as',
            '--with-gnu-ld',
        ],
        "powerpc-linux-gnu": [
            '--disable-sim',
            '--enable-lto',
            '--enable-relro',
            '--with-pic',
        ],
    }  # yapf: disable
    configure_arch_flags['aarch64-linux-gnu'] = [
        *configure_arch_flags['arm-linux-gnueabi'],
        '--enable-gold',
        '--enable-ld=default',
    ]
    configure_arch_flags['powerpc64-linux-gnu'] = configure_arch_flags[
        'powerpc-linux-gnu']
    configure_arch_flags['powerpc64le-linux-gnu'] = configure_arch_flags[
        'powerpc-linux-gnu']
    configure_arch_flags['riscv64-linux-gnu'] = configure_arch_flags[
        'powerpc-linux-gnu']
    configure_arch_flags['s390x-linux-gnu'] = [
        *configure_arch_flags['powerpc-linux-gnu'],
        '--enable-targets=s390-linux-gnu',
    ]
    configure_arch_flags['x86_64-linux-gnu'] = [
        *configure_arch_flags['powerpc-linux-gnu'],
        '--enable-targets=x86_64-pep',
    ]

    for endian in ["", "el"]:
        configure_arch_flags[f'mips{endian}-linux-gnu'] = [
            f'--enable-targets=mips64{endian}-linux-gnuabi64,mips64{endian}-linux-gnuabin32'
        ]

    configure += configure_arch_flags.get(target, [])

    # If the current machine is not the target, add the prefix to indicate
    # that it is a cross compiler
    if not host_is_target(target):
        configure += [f'--program-prefix={target}-', f'--target={target}']

    print_header(f"Building {target} binutils")
    subprocess.run(configure, check=True, cwd=build_folder)


def invoke_make(build_folder, install_folder, target):
    """
    Invoke make to compile binutils
    :param build_folder: Build directory
    :param install_folder: Directory to install binutils to
    :param target: Target to compile for
    """
    make = ['make', '-s', '-j' + str(multiprocessing.cpu_count()), 'V=0']
    if host_is_target(target):
        subprocess.run(make + ['configure-host'], check=True, cwd=build_folder)
    subprocess.run(make, check=True, cwd=build_folder)
    if install_folder:
        subprocess.run(make + [f'prefix={install_folder}', 'install'],
                       check=True,
                       cwd=build_folder)
        with install_folder.joinpath(".gitignore").open("w") as gitignore:
            gitignore.write("*")


def build_targets(binutils_folder, build, install_folder, targets, host_arch):
    """
    Builds binutils for all specified targets
    :param binutils_folder: Binutils source folder
    :param build: Build directory
    :param install_folder: Directory to install binutils to
    :param targets: Targets to compile binutils for
    :param host_arch: Host architecture to optimize for
    :return:
    """
    for target in targets:
        build_folder = build.joinpath(target)
        cleanup(build_folder)
        invoke_configure(binutils_folder, build_folder, install_folder, target,
                         host_arch)
        invoke_make(build_folder, install_folder, target)


def main():
    st = time.time()

    root_dir = pathlib.Path(__file__).resolve().parent

    args = parse_parameters(root_dir)

    if args.binutils_folder:
        binutils_folder = pathlib.Path(args.binutils_folder).resolve()
    else:
        binutils_folder = root_dir.joinpath(BINUTIL_VERSION)
        download_binutils(root_dir)

    build_folder = pathlib.Path(args.build_folder).resolve()

    if args.skip_install:
        install_folder = None
    else:
        install_folder = pathlib.Path(args.install_folder).resolve()

    targets = ["all"]
    if args.targets is not None:
        targets = args.targets

    build_targets(binutils_folder, build_folder, install_folder,
                  create_targets(targets), args.march)

    print(f"\nTotal script duration: {get_duration(st)}")


if __name__ == '__main__':
    main()

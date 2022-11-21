#include <stdio.h>

__attribute__((noinline))
void foo() { printf("foo\n"); }

int main(int argc, char **argv) {
    if (argc == 2)
        foo();
    printf("main\n");
}


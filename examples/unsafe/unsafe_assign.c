int unsafe_assign(int* buf){
	if (*buf > 1999){
		free(buf);
	}
	else {
		*buf = 99;
	}
	return 0;
}

int main(int argc, char **argv) {

    if (argc<2) {
        puts("./unsafe_assign <N>");
        return 1;
    }

	int* val = malloc(sizeof(int));
	*val = atoi(argv[2]);
    printf("Freeing %d\n", unsafe_assign(val));
    printf("%d\n", *val);

    return 0;
}

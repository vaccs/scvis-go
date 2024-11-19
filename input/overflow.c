#include <stdio.h>

extern void func(char* cptr);

int main(int argc, char *argv[]){
  printf("argv[0] = %s\n",*argv);
  printf("argv[1] = %s\n",argv[1]);
  func(argv[1]);
}

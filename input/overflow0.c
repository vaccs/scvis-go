#include <string.h>
#include <stdlib.h>
#include <stdio.h>

void overwrite(int x, char *buf){
  int local=5;
  char internal[8];

  printf("local is %d\n",local);  
  memcpy(internal,buf,x);
  printf("local is %d\n",local);
}
int main(int argc, char *argv[]){
  overwrite(atoi(argv[1]),argv[2]);
}

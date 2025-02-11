#include <string.h>
int i;
void func(char *cptr){
  char buf[8];
  int j;

  i=7;
  j=11;
  strcpy(buf,cptr);
  i=0;
  j=0;
}
int main(int argc, char *argv[]){
  func(argv[1]);
}

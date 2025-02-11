#include <stdio.h>
int a;
int c = 3;
int main() {

   int b;

   int *p = NULL;

   p = &b; 
   a = 1;
   b = 2;
   printf("a = %d, b = %d, *p = %d, c = %d\n",a, b, *p, c);
}

#include <stdio.h>
int a;
int c = 3;
int main() {

   float b;

   float *p = NULL;

   p = &b; 
   a = 1;
   b = 2.0;
   printf("a = %d, b = %f, *p = %f, c = %d\n",a, b, *p, c);
}

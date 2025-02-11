#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>
struct tests {
   int a;
   int b;
   int *p;
};
int main() {

   struct tests x;

   x.a = 1;
   x.b = 2;
   x.p = &x.b;
   printf("a = %d, b = %d, *p = %d\n",x.a, x.b,*x.p);
}

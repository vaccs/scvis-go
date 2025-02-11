#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>
int main() {
int a[10];

   for (int i = 0; i < 10; i++) {
      a[i] = i;
      printf("a[%d] = %d\n",i,a[i]);
   }
}

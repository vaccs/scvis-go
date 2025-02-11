#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int main() {

   char *a;

   a = malloc(8*sizeof(char));
   strcpy(a,"0123456789abcdef");   
   *a = 'v';
   printf("a = %s\n",a);
}

#include<stdio.h>
#include <sys/types.h>
#include <unistd.h>
#include<limits.h>
int a[10];
void bar(int c[]){
	int b[10];
        int i;
       printf("In bar()\n");
       for (i=0;i<10;i++){
          b[i] = i;
          c[i] = b[i] +1;
       }
}
void foo(){
   printf("In foo()\n");	
}

int main(){
	foo();
        bar(a);
}

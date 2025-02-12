#include <stdlib.h>
int  gy;
char *gdata;
void f(int parm)
{
  static int count=0;
  int f1 = 9;
  gy=f1+parm+count;
  *gdata=f1;
}
int main(){
  int m1;  
  int m2;  
  m1=3;
  m2=6;
  gdata=malloc(10);  
  f(m1+m2);

}  

#include <stdio.h>
#include <stdlib.h>
int afunc(int p1, short p2, char p3){
      int li;
      short ls;

      li=p1;
      ls=p3;
      li=li+p2+ls;
 
      return(li);
}
int main(int argc, char *argv[]){
  
      int mi;
      short ms;

      if (argc<4)exit(1);
      mi=atoi(argv[1]);
      ms=atoi(argv[2]);
      mi=afunc(mi,ms,argv[3][0]);

      printf("mi is %d\n",mi);
}      

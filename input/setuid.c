#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#define GNU_SOURCE
#include <unistd.h>
#include <sys/types.h>

int main() {
  uid_t rUid, eUid, sUid;

  getresuid(&rUid,&eUid,&sUid);
  printf("Real: <%d> Effective: <%d> Saved: <%d>\n",rUid,eUid,sUid);
  seteuid(rUid);
  getresuid(&rUid,&eUid,&sUid);
  printf("Real: <%d> Effective: <%d> Saved: <%d>\n",rUid,eUid,sUid);
  seteuid(sUid);
  getresuid(&rUid,&eUid,&sUid);
  printf("Real: <%d> Effective: <%d> Saved: <%d>\n",rUid,eUid,sUid);
}


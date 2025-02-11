#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>
#include <assert.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <stdarg.h>
#include <limits.h>

void vaccs_printf(char *control_string, ...) {

   va_list ap;
   static int stdout_fd = -1;
   static char *stdout_fn = "/tmp/vaccs.stdout";

   va_start (ap, control_string);
   if (stdout_fd == -1) {
      assert((stdout_fd = open(stdout_fn,O_WRONLY)) != -1);
   }

   char stdout_buff[UCHAR_MAX];
   int num_chars = vsnprintf(stdout_buff,UCHAR_MAX,control_string,ap);

   assert(num_chars >= 0 && num_chars < UCHAR_MAX);

   va_end(ap);

   printf("%s",stdout_buff);
   if (stdout_fd != -1)
      assert(write(stdout_fd,stdout_buff,num_chars) == num_chars);

}

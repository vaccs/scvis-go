int simple(int x, int y)
{
  int local;
  local=x+y;
  return(local);
}
void simple2(int x, int y)
{
  char c;
  c=x+y;
}
int main()
{
  int ret;
  ret=simple(7,13);
  simple2(-1,0);
  return(0);
}

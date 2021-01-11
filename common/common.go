package common

import (
  "os"
  "strconv"
)

func LookupStrEnv(name string, dflt string) (string) {
  strEnvStr, ok := os.LookupEnv(name)
  if ok == false {
    return dflt
  }

  return strEnvStr
}

func LookupIntEnv(name string, dflt int64) (int64) {
  intEnvStr, ok := os.LookupEnv(name)
  if ok == false {
    return dflt
  }

  intEnv, parseerr := strconv.ParseInt(intEnvStr, 10, 64)
  if parseerr != nil {
    return dflt
  }

  return intEnv
}

func LookupBooleanEnv(name string, dflt bool) (bool) {
  boolEnvStr, ok := os.LookupEnv(name)
  if ok == false {
    return dflt
  }
  boolEnv, parseerr := strconv.ParseBool(boolEnvStr)
  if parseerr != nil {
    return dflt
  }

  return boolEnv
}

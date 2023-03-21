type hookCall struct {
	method string
	params []interface{}
  }
  
  var (
	callch = make(chan hookCalls, 100)
  )
  var errs []errors
  
  func main() {
	select {
	case callRecord := <- callch:
	  //whatever
	case <-time.NewTimer(10 * time.second).C:
	  // err
	  
	}
  
	err := client.Call("my_method", &whatever, params...)
	// Check whatever is as expected and err == nil
	
	if callRecord.method != "getRPC" {
	  errs = append(errs, fmt.Errorf("Get RPC not called"))
	}
  
	if len(callRecord.params) != 2 {
	  // Report error
	}
	if _, ok := callRecord.params[1].(restricted.Params); !ok {
	  errs = append(errs, fmt.Errorf("params not of right type"))
	}
	
	// Use engine_newPlayloadV1 to insert a block
	callRecord = <-callch
	if callRecord.method != "XYZ" {
	  errs = append()
	}
  
	if len(errs) > 0 {
	  for _, err := range errs {
		log.Error("Error", "err", err)
	  }
	  os.Exit(1)
	}
	os.Exit(0)
  }
  
  func GetRPCMethods(x, y, z string) {
	callch <- hookCall{"getRPC", []interface{x, y, z}}
  }
  
  func ProcessBlock() {
	callch <- hookCall{"processBlock"}
  }
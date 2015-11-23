/*@Author :Manasvini Banavara Suryanarayana
*@SJSU ID : 010102040
*CMPE 273 Assignment #3
*/
package main

import (
    "fmt"
    "./httprouter"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "net/url"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "os"
    "strconv"
    "math"
    "strings"
    "bytes"
)

type Coordinate struct {
   Lat float64 `json:"lat"`
   Lng float64 `json:"lng"`
}

type ReqUberTrip struct{
    ProdId string `json:"product_id"`
    SLat float64 `json:"start_latitude"`
    SLong float64 `json:"start_longitude"`

} 

type Details struct {
   ID int32  `json:"id"`
   Name string `json:"name"`
   Address string `json:"address"`
   City string `json:"city"`
   State string `json:"state"`
   Zip string `json:"zip"`
   Coordinate Coordinate `json:"coordinate"`
 }
  
  type UpdDetails struct {
   Address string `json:"address"`
   City string `json:"city"`
   State string `json:"state"`
   Zip string `json:"zip"`
   Coordinate Coordinate `json:"coordinate"`
 }
  type RespBodyTrip struct{
    ID string `json:"id"`
    Status string `json:"status"`
    Startpoint string `json:"starting_from_location_id"`
    TripLocations []string `json:"best_route_location_ids"`
    TotalCost int32 `json:"total_uber_costs"`
    TotalDuration int32  `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"`

  }

  type RespBodyTripRequest struct{
    ID string `json:"id"`
    Status string `json:"status"`
    Startpoint string `json:"starting_from_location_id"`
    NextLoc string `json:"next_destination_location_id"`
    TripLocations []string `json:"best_route_location_ids"`
    TotalCost int32 `json:"total_uber_costs"`
    TotalDuration int32  `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"`
    WaitTimeETA int `json:"uber_wait_time_eta"`

  }

  type UberResp struct{
    Price float64 `json:"high_estimate"`
    Distance float64 `json:"duration"`
    Duration float64 `json:"distance"`

  }

func getmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
  var id string = p.ByName("id")
  id2, err9 :=strconv.Atoi(id)
  if err9 != nil {
        fmt.Println(err9)
    }
  fmt.Println(id)
    //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
 collection := sess.DB("cmpe273").C("table")


  var responsemesg Details

  err = collection.Find(bson.M{"id":id2}).One(&responsemesg)
  if err != nil {
    fmt.Printf("got an error finding a doc %v\n",err)
    os.Exit(1)
  }

    //converting response body struct to json format
   respjson, err5 := json.Marshal(responsemesg)
   if err5 != nil {
        fmt.Println(err5)
    }
     
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(200)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)
    
}

func postmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
   fmt.Println("1")
   //creating struct to read request json 
   type ReqBody struct {
        Name string `json:"name"`
        Address string `json:"address"`
        City string `json:"city"`
        State string `json:"state"`
        Zip string `json:"zip"`
    }  
    var x ReqBody

    //fetch the body from request 
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
       fmt.Println(err)
   }
    
    //converting json body to struct of type ReqBody
    err1 := json.Unmarshal(body, &x)
    if err1 != nil {
       fmt.Println(err1)
   }
  //creating google api query
   var x1 string = x.Address+","+x.City+","+x.State
    fmt.Println("Query sent to google API : " + x1)
    var input string = url.QueryEscape(x1)
    resp, err2 := http.Get("http://maps.google.com/maps/api/geocode/json?address="+input+"&sensor=false")

    if err2 != nil {
        fmt.Println(err2)
    }
    defer resp.Body.Close()

    body1, err3 := ioutil.ReadAll(resp.Body)
    if err3 != nil {
        fmt.Println(err3)
    }

    var googmap interface{}
    err4 := json.Unmarshal(body1, &googmap)
    if err4 != nil {
        fmt.Println(err4)
    }

    res1 := googmap.(map[string]interface{})
    i := res1["results"]
    m1 := i.([]interface{})
    x8 := m1[0]
    x2 := x8.(map[string]interface{})
     y1 := x2["geometry"]
     y2 := y1.(map[string]interface{})
     y3 := y2["location"]
     fmt.Println(y3)
     y4 := y3.(map[string]interface{})
     latval := y4["lat"]
     lngval := y4["lng"]
     
   
     fmt.Println("3")
    //constructing struct for sending back response body
     cord := Coordinate{
      Lat: latval.(float64),
      Lng: lngval.(float64),
     }
    fmt.Println("4")
    var idval int32 = bson.NewObjectId().Counter()
    fmt.Println(idval)
    test := Details{
        ID:  idval,
        Name: x.Name,
        Address: x.Address,
        City: x.City,
        State: x.State,
        Zip: x.Zip,
        Coordinate: cord,
      }

      //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table")
 
  err = collection.Insert(test)
  if err != nil {
    fmt.Printf("Can't insert document: %v\n", err)
    os.Exit(1)
  }

    //converting response body struct to json format
   respjson, err5 := json.Marshal(test)
   if err5 != nil {
        fmt.Println(err5)
    }
     
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(201)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)
     
}
func putmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    //fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("id"))
  //fetching id from request
  var id string = p.ByName("id")
  id2, err9 :=strconv.ParseInt(id,10,32)

  if err9 != nil {
        fmt.Println(err9)
    }

    //fetching data from request body

    //creating struct to read request json 
   type ReqBody struct {
        Address string `json:"address"`
        City string `json:"city"`
        State string `json:"state"`
        Zip string `json:"zip"`
    }  
    var x ReqBody

    //fetch the body from request 
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
       fmt.Println(err)
   }
    
    //converting json body to struct of type ReqBody
    err1 := json.Unmarshal(body, &x)
    if err1 != nil {
       fmt.Println(err1)
   }

   //creating google api query
   var x1 string = x.Address+","+x.City+","+x.State
    fmt.Println("Query sent to google API : " + x1)
    var input string = url.QueryEscape(x1)
    resp, err2 := http.Get("http://maps.google.com/maps/api/geocode/json?address="+input+"&sensor=false")

    if err2 != nil {
        fmt.Println(err2)
    }
    defer resp.Body.Close()

    body1, err3 := ioutil.ReadAll(resp.Body)
    if err3 != nil {
        fmt.Println(err3)
    }

    var googmap interface{}
    err4 := json.Unmarshal(body1, &googmap)
    if err4 != nil {
        fmt.Println(err4)
    }

    res1 := googmap.(map[string]interface{})
    i := res1["results"]
    m1 := i.([]interface{})
    x8 := m1[0]
    x2 := x8.(map[string]interface{})
     y1 := x2["geometry"]
     y2 := y1.(map[string]interface{})
     y3 := y2["location"]
     fmt.Println(y3)
     y4 := y3.(map[string]interface{})
     latval := y4["lat"]
     lngval := y4["lng"]
     

     //constructing struct for sending back response body
     cord := Coordinate{
      Lat: latval.(float64),
      Lng: lngval.(float64),
     }
    fmt.Println("4")
    
    test := UpdDetails{
        Address: x.Address,
        City: x.City,
        State: x.State,
        Zip: x.Zip,
        Coordinate: cord,
      }

    //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table")

  //update the records for particular ID
  err10 := collection.Update(bson.M{"id":id2},bson.M{"$set": test})
  if err10 != nil {
    fmt.Printf("Can't update document: %v\n", err10)
    os.Exit(1)
  }

  //fetch the record to get name
  var responsemesg Details

  err11 := collection.Find(bson.M{"id":id2}).One(&responsemesg)
  if err11 != nil {
    fmt.Printf("got an error finding a doc %v\n",err11)
    os.Exit(1)
  }

    //converting response body struct to json format
   respjson, err5 := json.Marshal(responsemesg)
   if err5 != nil {
        fmt.Println(err5)
    }
     
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(201)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)

}
func delmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    

  //fetching id from request
  var id string = p.ByName("id")
  id2, err9 :=strconv.ParseInt(id,10,32)

  if err9 != nil {
        fmt.Println(err9)
    }

  //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table")

  //delete the particular id record

  err10 := collection.Remove(bson.M{"id":id2})
  if err10 != nil {
    fmt.Printf("Can't Remove document: %v\n", err10)
    os.Exit(1)
  }

  rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(200)

}
//this method fetche latitude/longitude of the locations ID from dm
func getLatLang(locId string) Details {

  id2, err9 :=strconv.Atoi(locId)
  if err9 != nil {
        fmt.Println(err9)
    }

   fmt.Println("Getting Lang : --|"+locId+"|---")
//database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
 collection := sess.DB("cmpe273").C("table")


  var responsemesg Details

  err = collection.Find(bson.M{"id":id2}).One(&responsemesg)
  if err != nil {
    fmt.Printf("got an error finding a doc %v\n",err)
    os.Exit(1)
  }

  return responsemesg

} 

//method returns uber trip response
func uberResponse(startPoint Details,locationSet []Details)RespBodyTrip{
OrderedLocID := make([]Details, len(locationSet))
LocationIdSet := make([]int32, len(locationSet))
var  MinDuration float64
 MinDuration = math.MaxFloat64
 var MinCost  float64
 MinCost = math.MaxFloat64
 MinDist := math.MaxFloat64
 ProdVar := math.MaxFloat64
 var SumDuration float64
 var SumCost float64
 SumDist := 0.0
 var curloc Details
 Startpt := startPoint

 for i :=0;i<len(locationSet);i++{
  fmt.Println(" i ",i)
  fmt.Println("start ",Startpt.ID)
  for j,s := range locationSet{
    flag := true
    fmt.Println("j ",j)
    fmt.Println("s : ",s.ID)
    for k,s1 := range LocationIdSet{
      fmt.Println("k ",k)
      fmt.Println("s1 - ",s1)
      if s1 == s.ID{
        flag=false
        fmt.Println("Skip")
        break
      }
    }
    if flag{
      fmt.Println("Callling Uber")
      fmt.Println(Startpt.ID)
      fmt.Println(s.ID)
      UberResult := callUber(Startpt,s)
      temp := ((UberResult.Price * UberResult.Duration)/10)
      if temp < ProdVar{
        MinDist=UberResult.Distance
        MinDuration = UberResult.Duration
        MinCost =UberResult.Price
        curloc = s
        ProdVar=temp
      }
    }
    
  }
  Startpt = curloc
  LocationIdSet[i]=curloc.ID
  fmt.Println("Location Id : ", curloc.ID)
  OrderedLocID[i] = curloc
  SumDuration = SumDuration + MinDuration
  SumCost =SumCost + MinCost
  SumDist =SumDist + MinDist
  
  ProdVar = math.MaxFloat64

 }
//for responsetrip body object to return as final response
 UberResult2 := callUber(curloc,startPoint)
 SumDuration = SumDuration + UberResult2.Duration
  SumCost =SumCost + UberResult2.Price
  SumDist =SumDist + UberResult2.Distance
 //converting Location id set from Int32 to string
 LocationIdSet2 := make([]string, len(LocationIdSet))
  for m,l := range LocationIdSet{
    LocationIdSet2[m]=strconv.Itoa(int(l))
 }
 var tripId int = int(bson.NewObjectId().Counter())
 temp := (tripId/1000)
  RetResp := RespBodyTrip{
    ID: strconv.Itoa(temp),
    Status: "planning",
    Startpoint: strconv.Itoa(int(startPoint.ID)),
    TripLocations: LocationIdSet2,
    TotalCost: int32(SumCost),
    TotalDuration: int32(SumDuration),
    TotalDistance: SumDist,
  }

return RetResp
}
//calling uber app to fetch price set of json
func callUber(sp Details, dp Details) UberResp{
  var x1 string = "start_latitude="+strconv.FormatFloat(sp.Coordinate.Lat, 'f', 7, 64)+"&start_longitude="+strconv.FormatFloat(sp.Coordinate.Lng, 'f', 7, 64)+"&end_latitude="+strconv.FormatFloat(dp.Coordinate.Lat, 'f', 7, 64)+"&end_longitude="+strconv.FormatFloat(dp.Coordinate.Lng, 'f', 7, 64)
    //var x1 string = "start_latitude=37.3976151&start_longitude=-122.0926371&end_latitude=37.3348765&end_longitude=-121.8887361"
  //var x1 string = "start_latitude=&start_longitude=&end_latitude=&end_longitude="
    fmt.Println("Query sent to uber API : " + x1)
    //var input string = url.QueryEscape(x1)
    client := &http.Client{}
    req, err := http.NewRequest("GET",  "https://api.uber.com/v1/estimates/price?"+x1, nil)
    if err != nil {
        fmt.Println(err)
    }
    req.Header.Set("Authorization", "Token -Xn1IwcO856_ow-sN-_jzkzkNQlc5OqDVV3w1aUR")
    resp, err2 := client.Do(req)
    //resp, err2 := http.Get("https://api.uber.com/v1/price?"+input)

    if err2 != nil {
        fmt.Println(err2)
    }
    defer resp.Body.Close()

    body1, err3 := ioutil.ReadAll(resp.Body)
    if err3 != nil {
        fmt.Println(err3)
    }
    //fmt.Println(body1)
    var uberapi interface{}
    err4 := json.Unmarshal(body1, &uberapi)
    if err4 != nil {
        fmt.Println(err4)
    }
    fmt.Println(uberapi)
    res1 := uberapi.(map[string]interface{})
    i := res1["prices"]
    m1 := i.([]interface{})
    x8 := m1[0]
    x2 := x8.(map[string]interface{})
     est := x2["high_estimate"]
    dur := x2["duration"]
     dist := x2["distance"]
     fmt.Println(est)
     fmt.Println(dur)
     fmt.Println(dist)
     resp1 := UberResp{
      Price: est.(float64),
      Distance: dist.(float64),
      Duration: dur.(float64),
     }
     return resp1
}
//adding trip method 

func tripmethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
   fmt.Println("tripmethod inside")
   //creating struct to read request json 
   type ReqTripBody struct {
        StartLocID string `json:"starting_from_location_id"`
        TripLocations []string `json:"location_ids"`
        
    }  
    var x ReqTripBody


  //fetch the body from request 
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
       fmt.Println(err)
   }
    
    //converting json body to struct of type ReqBody
    err1 := json.Unmarshal(body, &x)
    if err1 != nil {
       fmt.Println(err1)
   }
   // calling getLatLang method to fetch co ordinates
  startPoint := getLatLang(x.StartLocID)
  locs := make([]Details, len(x.TripLocations))
  for i,s := range x.TripLocations {
    fmt.Println(i,s)
    locs[i]=getLatLang(s)


  }
  SendResp := uberResponse(startPoint,locs)

   //database connection
    sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
    }
    defer sess.Close()

  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table_trip")
 
  err = collection.Insert(SendResp)
  if err != nil {
    fmt.Printf("Can't insert document: %v\n", err)
    os.Exit(1)
  }


   //converting response body struct to json format
   respjson, err5 := json.Marshal(SendResp)
   if err5 != nil {
        fmt.Println(err5)
    }
    
    fmt.Println("sent")
   rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(201)
    //sending back response
    fmt.Fprintf(rw, "%s", respjson)
     
}
func tripGetMethod(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
  var responsemesg1 RespBodyTrip
  var reqId string = p.ByName("id")
  fmt.Println("---|"+reqId+"|----")
  sess1, err8 := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
  if err8 != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err8)
    os.Exit(1)
  }
  defer sess1.Close()
  sess1.SetSafe(&mgo.Safe{})
  collection3 := sess1.DB("cmpe273").C("table_trip")
  err2 := collection3.Find(bson.M{"id":reqId}).One(&responsemesg1)
  if err2 != nil {
    fmt.Printf("got an error finding a doc %v\n",err2)
    os.Exit(1)
  }
  respjson, err5 := json.Marshal(responsemesg1)
  if err5 != nil {
    fmt.Println(err5)
  }
    
  fmt.Println("sent")
  rw.Header().Set("Content-Type","application/json")
  rw.WriteHeader(200)
  //sending back response
  fmt.Fprintf(rw, "%s", respjson)
}

//trip request method begins here 
func tripRequest(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
  var reqId string = p.ByName("id")
  id2, err9 :=strconv.Atoi(reqId)
  if err9 != nil {
    fmt.Println(err9)
  }
  fmt.Println(id2)

  //getting auth token from header
  var authtoken string
  
  if req.Header.Get("Authorization")=="" {
    fmt.Println("Invalid auth Token")
   
    Noauthmesg := "message:Auth token is missing please include it in the header of the request in the specified format, format: Authorization: Bearer <auth_token>"
    
    rw.WriteHeader(401)
    //sending back response
    fmt.Fprintf(rw, "%s", Noauthmesg)
    return
  }

  fmt.Println("Authorization")
  authtoken = strings.Split(req.Header.Get("Authorization")," ")[1]
  fmt.Println(authtoken)


  //database connection
  sess, err := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
  if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
  }
  defer sess.Close()
  sess.SetSafe(&mgo.Safe{})
  collection := sess.DB("cmpe273").C("table_triprequest")
  var responsemesg RespBodyTripRequest
  var responsemesg1 RespBodyTrip
  var Details1 Details
  flag := 0
  err = collection.Find(bson.M{"id":reqId}).One(&responsemesg)
  if err != nil {
    sess1, err8 := mgo.Dial("mongodb://admin:admin@ds043714.mongolab.com:43714/cmpe273")
    if err8 != nil {
      fmt.Printf("Can't connect to mongo, go error %v\n", err8)
      os.Exit(1)
    }
    defer sess1.Close()
    sess1.SetSafe(&mgo.Safe{})
    fmt.Println("1st Attempt")
    collection3 := sess1.DB("cmpe273").C("table_trip")
    err2 := collection3.Find(bson.M{"id":reqId}).One(&responsemesg1)
    if err2 != nil {
      fmt.Printf("got an error finding a doc %v\n",err2)
      os.Exit(1)
    }
    startpoint1 := responsemesg1.Startpoint
    Details1 = getLatLang(startpoint1)
    flag =1
    fmt.Println("Found in 2nd table")
  } else {
    fmt.Println("Found in 3rd table")
    if responsemesg.Status == "finished" {
      respjson, err5 := json.Marshal(responsemesg)
      if err5 != nil {
        fmt.Println(err5)
      }
      var responsemesg2 RespBodyTrip
      collection3 := sess.DB("cmpe273").C("table_trip")
      err2 := collection3.Find(bson.M{"id":reqId}).One(&responsemesg2)
      if err2 != nil {
        fmt.Printf("got an error finding a doc %v\n",err2)
        os.Exit(1)
      }
      responsemesg2.Status = "finished"
      //update the records for particular ID
      err10 := collection3.Update(bson.M{"id":reqId},bson.M{"$set": responsemesg2})
      if err10 != nil {
        fmt.Printf("Can't update document: %v\n", err10)
        os.Exit(1)
      }
      fmt.Println("sent")
      rw.Header().Set("Content-Type","application/json")
      rw.WriteHeader(201)
      //sending back response
      fmt.Fprintf(rw, "%s", respjson)
      return 
    }
    startpoint := responsemesg.NextLoc    
    Details1 = getLatLang(startpoint) 
  }
  fmt.Println(Details1.ID)
  fmt.Println("Call Uber")
  // call uber api and get ETA

  group := ReqUberTrip {
    ProdId: "04a497f5-380d-47f2-bf1b-ad4cfdcb51f2",
    SLat: Details1.Coordinate.Lat,
    SLong: Details1.Coordinate.Lng,
  }
  b, err5 := json.Marshal(group)
  if err5 != nil {
    fmt.Println(err5)
  } 

  client1 := &http.Client{}
  req, err111 := http.NewRequest("POST",  "https://sandbox-api.uber.com/v1/requests", bytes.NewBuffer(b))
  if err111 != nil {
    fmt.Println(err111)
  }
  req.Header.Set("Authorization", "Bearer "+authtoken)
  req.Header.Set("Content-Type", "application/json")
  resp, err2 := client1.Do(req)
  
  if err2 != nil {
    fmt.Println(err2)
  }
  defer resp.Body.Close()
  body1, err3 := ioutil.ReadAll(resp.Body)
  if err3 != nil {
    fmt.Println(err3)
  }
  var uberapi interface{}
  err4 := json.Unmarshal(body1, &uberapi)
  if err4 != nil {
    fmt.Println(err4)
  }
  fmt.Println(uberapi)
  res1 := uberapi.(map[string]interface{})
  i := res1["eta"]     
  eta := int(i.(float64))
  fmt.Println(int(i.(float64)))

  
  //constructuing response body
  var respbody RespBodyTripRequest
  if flag == 0 {
    var nextLocation string
    status := responsemesg.Status
    // Calculate next location
    for i,s := range responsemesg.TripLocations{
      fmt.Println(i,s)
      if s == responsemesg.NextLoc{
        if i == (len(responsemesg.TripLocations)-1){
          status = "finished"
          nextLocation = responsemesg.NextLoc
        } else {
          index := (i+1)
          nextLocation = responsemesg.TripLocations[index]
        }
        break; 
      }
    }
    respbody = RespBodyTripRequest{
      ID: responsemesg.ID,
      Status: status,
      Startpoint: responsemesg.Startpoint,
      NextLoc: nextLocation,
      TripLocations: responsemesg.TripLocations,
      TotalCost: responsemesg.TotalCost,
      TotalDuration: responsemesg.TotalDuration,
      TotalDistance: responsemesg.TotalDistance,
      WaitTimeETA: eta,
    }
    collection2 := sess.DB("cmpe273").C("table_triprequest")
    //update the records for particular ID
    err10 := collection2.Update(bson.M{"id":reqId},bson.M{"$set": respbody})
    if err10 != nil {
      fmt.Printf("Can't update document: %v\n", err10)
      os.Exit(1)
    }
  } else {
    respbody = RespBodyTripRequest{
      ID: responsemesg1.ID,
      Status: "requesting",
      Startpoint: responsemesg1.Startpoint,
      NextLoc: responsemesg1.TripLocations[0],
      TripLocations: responsemesg1.TripLocations,
      TotalCost: responsemesg1.TotalCost,
      TotalDuration: responsemesg1.TotalDuration,
      TotalDistance: responsemesg1.TotalDistance,
      WaitTimeETA: eta,
    }
    responsemesg1.Status = "requesting"
    collection5 := sess.DB("cmpe273").C("table_trip")
    //update the records for particular ID
    err10 := collection5.Update(bson.M{"id":reqId},bson.M{"$set": responsemesg1})
    if err10 != nil {
      fmt.Printf("Can't update document: %v\n", err10)
      os.Exit(1)
    }
    collection4 := sess.DB("cmpe273").C("table_triprequest")
    err = collection4.Insert(respbody)
    if err != nil {
      fmt.Printf("Can't insert document: %v\n", err)
      os.Exit(1)
    }  
  }
  

  //converting response body struct to json format
  respjson, err5 := json.Marshal(respbody)
  if err5 != nil {
    fmt.Println(err5)
  }
    
  fmt.Println("sent")
  rw.Header().Set("Content-Type","application/json")
  rw.WriteHeader(201)
  //sending back response
  fmt.Fprintf(rw, "%s", respjson)

}

func main() {
    mux := httprouter.New()
    mux.GET("/locations/:id", getmethod)
    mux.POST("/locations", postmethod)
    mux.PUT("/locations/:id", putmethod)
    mux.DELETE("/locations/:id", delmethod)
    mux.POST("/trips",tripmethod)
    mux.PUT("/trips/:id/request",tripRequest)
    mux.GET("/trips/:id", tripGetMethod)
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }
    fmt.Println("serverstarted")
    server.ListenAndServe()
}
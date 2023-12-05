# gin-params-mapper

gin-params-mapper provide functions to find and convert params from a `gin.Context` object.

You can get request param from the Query, the Path or the Context itself.

For all functions you have to pass the param `name`, the `gin.Context` and the **address** of the destination object.

## Supported destination types :
- string
- float64
- int
- int64
- bool
- []byte
- uuid.UUID : the param must be a string in the go UUID format
- time.Time : the param must be a string in the RFC 3339 Nano format
- time.Duration : the param must be a string in the go format (exemple : 1s)

## Exemple :

```golang
func GetFoo(c *gin.Context) {
    var bar uuid.UUID

    if err := ginparamsmapper.GetFromContext(ginparamsmapper.QueryLocation ,"bar", c, &bar); err != nil {
        c.JSON(400, "Can't find the param bar from the query. Error : "+err.Error())
        return
    }

    log.Debug(bar)
}
```

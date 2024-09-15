package api

// func AddRoutes(r *gin.Engine) {
// 	r.GET("/api/v1/contexts", func(c *gin.Context) {
// 		properties := packinglib.ListProperties()
// 		propDescriptions := make([][2]string, len(properties))
// 		for i, p := range properties {
// 			propDescriptions[i] = [2]string{string(p), p.Description()}
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"contexts": propDescriptions,
// 		})
// 	})
// }

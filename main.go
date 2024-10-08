package main

import (
    "log"
    "net/http"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
    "github.com/labstack/echo/v5"
)

func msalMiddleware(app core.App) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            return next(c) //TODO remove this for middleware
            // Your MSAL token verification logic here
            token := c.Request().Header.Get("Authorization")
            if token == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No token provided"})
            }
            // Add your MSAL token verification logic here
            // If verification fails, return an error response
            // For now, we're just checking if the token exists
            return next(c)
        }
    }
}

func main() {
    app := pocketbase.New()

    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        // Apply MSAL middleware to all routes
        e.Router.Use(msalMiddleware(app))

        // Specifically for the users collection
        e.Router.AddRoute(echo.Route{
            Method: http.MethodGet,
            Path:   "/api/collections/users/records",
            /* Handler: func(c echo.Context) error {
                return apis.RequireAdminOrRecordAuth(apis.CollectionRecordsList)(c, app)
            }, */
            Middlewares: []echo.MiddlewareFunc{
                apis.RequireAdminOrRecordAuth(),
            },
        })

        e.Router.AddRoute(echo.Route{
            Method: http.MethodPost,
            Path:   "/api/collections/users/records",
            /* Handler: func(c echo.Context) error {
                return apis.RequireAdminOrRecordAuth(apis.CollectionRecordsCreate)(c, app)
            }, */
            Middlewares: []echo.MiddlewareFunc{
                apis.RequireAdminOrRecordAuth(),
            },
        })

        e.Router.AddRoute(echo.Route{
            Method: http.MethodPatch,
            Path:   "/api/collections/users/records/:id",
            /* Handler: func(c echo.Context) error {
                return apis.RequireAdminOrRecordAuth(apis.CollectionRecordsUpdate)(c, app)
            }, */
            Middlewares: []echo.MiddlewareFunc{
                apis.RequireAdminOrRecordAuth(),
            },
        })

        e.Router.AddRoute(echo.Route{
            Method: http.MethodDelete,
            Path:   "/api/collections/users/records/:id",
            /* Handler: func(c echo.Context) error {
                return apis.RequireAdminOrRecordAuth(apis.CollectionRecordsDelete)(c, app)
            }, */
            Middlewares: []echo.MiddlewareFunc{
                apis.RequireAdminOrRecordAuth(),
            },
        })

        // Add a custom route as an example
        e.Router.GET("/api/custom_user_endpoint", func(c echo.Context) error {
            return c.JSON(200, map[string]string{"message": "This is a custom user endpoint"})
        })

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}

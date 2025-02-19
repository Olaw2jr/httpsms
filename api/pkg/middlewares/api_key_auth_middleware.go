package middlewares

import (
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// APIKeyAuth authenticates a user from the X-API-Key header
func APIKeyAuth(logger telemetry.Logger, tracer telemetry.Tracer, userRepository repositories.UserRepository) fiber.Handler {
	logger = logger.WithService("middlewares.APIKeyAuth")

	return func(c *fiber.Ctx) error {
		ctx, span := tracer.StartFromFiberCtx(c, "middlewares.APIKeyAuth")
		defer span.End()

		ctxLogger := tracer.CtxLogger(logger, span)

		apiKey := c.Get(authHeaderAPIKey)
		if len(apiKey) == 0 || apiKey == "undefined" {
			span.AddEvent(fmt.Sprintf("the request header has no [%s] header", authHeaderAPIKey))
			return c.Next()
		}

		authUser, err := userRepository.LoadAuthUser(ctx, apiKey)
		if err != nil {
			ctxLogger.Error(stacktrace.Propagate(err, fmt.Sprintf("cannot load user with api key [%s]", apiKey)))
			return c.Next()
		}
		c.Locals(ContextKeyAuthUserID, authUser)
		ctxLogger.Info(fmt.Sprintf("[%T] set successfully for user with ID [%s]", authUser, authUser.ID))
		return c.Next()
	}
}

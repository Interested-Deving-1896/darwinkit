package usernotifications

import (
	"testing"

	"github.com/progrium/darwinkit/internal/assert"
)
	
func TestUserNotificationsValid(t *testing.T) {
	// Test that the classes can be accessed
	assert.NotNil(t, NotificationCenterClass)
	assert.NotNil(t, NotificationRequestClass)
	assert.NotNil(t, NotificationContentClass)
	assert.NotNil(t, MutableNotificationContentClass)
	assert.NotNil(t, NotificationTriggerClass)
	assert.NotNil(t, TimeIntervalNotificationTriggerClass)
	assert.NotNil(t, CalendarNotificationTriggerClass)
	assert.NotNil(t, LocationNotificationTriggerClass)
	assert.NotNil(t, PushNotificationTriggerClass)
	assert.NotNil(t, NotificationSoundClass)

	// Test the constants
	assert.Equal(t, AuthorizationOptionBadge, AuthorizationOptions(1))
	assert.Equal(t, AuthorizationOptionSound, AuthorizationOptions(2))
	assert.Equal(t, AuthorizationOptionAlert, AuthorizationOptions(4))
	assert.Equal(t, AuthorizationOptionCarPlay, AuthorizationOptions(8))
}
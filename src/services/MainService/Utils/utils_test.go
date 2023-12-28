package Utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsernameValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected int
	}{
		{
			name:     "valid email",
			username: "alifakhary622@gmail.com",
			expected: 0,
		},
		{
			name:     "valid phone number",
			username: "+989123456789",
			expected: 2,
		},
		{
			name:     "valid phone number",
			username: "09123456789",
			expected: 4,
		},
		{
			name:     "invalid username",
			username: "invalid_username",
			expected: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := UsernameValidation(test.username)
			if actual != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestEmailRandomGenerator(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "valid email",
			expected: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := EmailRandomGenerator()
			if actual == "" {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		})
	}
}

const (
	testConnectionString = "amqp://guest:guest@rabbitmq:5672/"
	testExchangeName     = "test_exchange"
	testExchangeType     = "direct"
	testQueueName        = "test_queue"
)

func TestConnectBroker_Success(t *testing.T) {
	client := MessageClient{}

	err := client.ConnectBroker(testConnectionString)

	assert.NoError(t, err, "Expected no error connecting to broker")
	assert.NotNil(t, client.Connection, "Expected a non-nil connection")
	client.Close()
}

func TestConnectBroker_EmptyConnectionString(t *testing.T) {
	client := MessageClient{}

	err := client.ConnectBroker("")

	assert.Error(t, err, "Expected an error for empty connection string")
	assert.Nil(t, client.Connection, "Expected a nil connection")
}

func TestPublishToExchange_Success(t *testing.T) {
	client := MessageClient{}
	err := client.ConnectBroker(testConnectionString)
	assert.NoError(t, err, "Expected no error connecting to broker")
	defer client.Close()

	message := []byte("Hello, Exchange!")

	err = client.Publish(message, testExchangeName, testExchangeType)

	assert.NoError(t, err, "Expected no error publishing to exchange")
}

func TestPublishToQueue_Success(t *testing.T) {
	client := MessageClient{}
	err := client.ConnectBroker(testConnectionString)
	assert.NoError(t, err, "Expected no error connecting to broker")
	defer client.Close()

	message := []byte("Hello, Queue!")

	err = client.PublishOnQueue(message, testQueueName)

	assert.NoError(t, err, "Expected no error publishing to queue")
}

func TestPublishToQueue_InvalidConnection(t *testing.T) {
	client := MessageClient{}

	message := []byte("Hello, Queue!")

	err := client.PublishOnQueue(message, testQueueName)

	assert.Error(t, err, "Expected an error publishing to queue with invalid connection")
}

const (
	testRedisHost     = "redis"
	testRedisPort     = "6379"
	testRedisPassword = ""
	testRedisDb       = 0
)

func TestConnectToRedis_Success(t *testing.T) {
	redisClient := NewRedisClient(testRedisHost, testRedisPort, testRedisPassword, testRedisDb)

	client := redisClient.Connect()
	defer client.Close()

	pong, err := client.Ping(redisClient.Context).Result()

	assert.NoError(t, err, "Expected no error connecting to Redis")
	assert.Equal(t, "PONG", pong, "Expected PONG response from Redis")
}

func TestSetAndGetRedisValue_Success(t *testing.T) {
	redisClient := NewRedisClient(testRedisHost, testRedisPort, testRedisPassword, testRedisDb)

	key := "test_key"
	value := "test_value"

	err := redisClient.Set(key, value)
	assert.NoError(t, err, "Expected no error setting value in Redis")

	result, err := redisClient.Get(key)
	assert.NoError(t, err, "Expected no error getting value from Redis")
	assert.Equal(t, value, result, "Expected retrieved value to match the set value")
}

func TestDeleteRedisKey_Success(t *testing.T) {
	redisClient := NewRedisClient(testRedisHost, testRedisPort, testRedisPassword, testRedisDb)

	key := "test_key"
	value := "test_value"

	err := redisClient.Set(key, value)
	assert.NoError(t, err, "Expected no error setting value in Redis")

	err = redisClient.Del(key)
	assert.NoError(t, err, "Expected no error deleting key in Redis")

	result, err := redisClient.Get(key)
	assert.Error(t, err, "Expected an error getting value for deleted key")
	assert.Empty(t, result, "Expected empty result for deleted key")
}

package emailservices_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/emailservices"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/stretchr/testify/require"
)

func TestAddToQueueMemory(t *testing.T) {
	servicer := emailservices.MakeMemoryMailer(logging.NewNopLogger())
	onCompletedCalled := false
	toTarget := "Harry.Potter@hogwarts.edu"

	firstEmailJob := emailservices.EmailJob{
		To:      toTarget,
		From:    "Hogwarts Admin<albus.dumbledore@hogwarts.edu>",
		Subject: "Tuition",
		Body:    "Gimme my money!",
		OnCompleted: func(err error) {
			onCompletedCalled = true
			require.NoError(t, err)
		},
	}

	servicer.AddToQueue(firstEmailJob)

	actualJob := servicer.Outbox[toTarget][0]

	require.Equal(t, firstEmailJob.Body, actualJob.Body)
	require.Equal(t, firstEmailJob.To, actualJob.To)
	require.Equal(t, firstEmailJob.From, actualJob.From)
	require.Equal(t, firstEmailJob.Subject, actualJob.Subject)
	require.True(t, onCompletedCalled)

	secondEmailJob := emailservices.EmailJob{
		To:      toTarget,
		From:    "royal.prince.",
		Subject: "Your kind Attention: Please help me move my vast fortune",
		Body: "Hello,\n\nAs you are aware, my father, the king of Nigeria, has been deposed and " +
			"killed. Before fleeing, I was able to transfer a large amount of money out of the country before this happened, " +
			"and left this money for you to claim, so long as you have my password. Kindly reply " +
			"so that we can discuss how to transfer this money to you.",
		OnCompleted: func(err error) {
			onCompletedCalled = true
			require.NoError(t, err)
		},
	}

	onCompletedCalled = false
	servicer.AddToQueue(secondEmailJob)
	actualJob = servicer.Outbox[toTarget][1]

	require.Equal(t, secondEmailJob.Body, actualJob.Body)
	require.Equal(t, secondEmailJob.To, actualJob.To)
	require.Equal(t, secondEmailJob.From, actualJob.From)
	require.Equal(t, secondEmailJob.Subject, actualJob.Subject)
	require.True(t, onCompletedCalled)

	fromTarget := "Hogwarts Admin<albus.dumbledore@hogwarts.edu>"
	thirdEmailJob := emailservices.EmailJob{
		To:      fromTarget,
		From:    toTarget,
		Subject: "Re: Tuition",
		Body:    "Do you accept Nigerian Naira?",
		OnCompleted: func(err error) {
			onCompletedCalled = true
			require.NoError(t, err)
		},
	}

	onCompletedCalled = false
	servicer.AddToQueue(thirdEmailJob)
	actualJob = servicer.Outbox[fromTarget][0] // since this is a reply, use from

	require.Equal(t, thirdEmailJob.Body, actualJob.Body)
	require.Equal(t, thirdEmailJob.To, actualJob.To)
	require.Equal(t, thirdEmailJob.From, actualJob.From)
	require.Equal(t, thirdEmailJob.Subject, actualJob.Subject)
	require.True(t, onCompletedCalled)
}

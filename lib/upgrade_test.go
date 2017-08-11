package cookie_clicker

import (
	"fmt"
	"testing"
	"time"
)

const BUILDING_UPGRADE_NAME string = "upgrade name"
const BUILDING_UPGRADE_RATIO float64 = 2
const BUILDING_UPGRADE_COST float64 = 100

func TestMakeBuildingUpgrade(t *testing.T) {
	var b BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)
	if b.name != BUILDING_UPGRADE_NAME {
		t.Error(fmt.Sprintf("Expected name \"%s\", got \"%s\"", BUILDING_UPGRADE_NAME, b.name))
	}
	if b.upgrade_ratio != BUILDING_UPGRADE_RATIO {
		t.Error(fmt.Sprintf("Expected ratio %e, got %e", BUILDING_UPGRADE_RATIO, b.upgrade_ratio))
	}
	if b.cost != BUILDING_UPGRADE_COST {
		t.Error(fmt.Sprintf("Expected cost %e, got %e", BUILDING_UPGRADE_COST, b.cost))
	}
}

func TestActivateBuildingUpgradeLoop(t *testing.T) {
	var b BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)

	go activateBuildingUpgradeLoop(&b)

	<-b.activated_send_channel

	select {
	case <-b.activated_send_channel:
		t.Error("Got unexpected building upgrade activated send signal")
	case <-time.After(CHANNEL_TIMEOUT):
	}
}

func TestActivateBuildingUpgrade(t *testing.T) {
	var b BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)

	go activateBuildingUpgrade(&b)

	for i := 0; i < 10; i++ {
		if !<-b.activated_send_channel {
			t.Error("Expected channel activated, was not")
		}
	}

	b.activated_send_done_channel <- true
}

func TestListenBuildingUpgradeLoop(t *testing.T) {
	var b BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)

	go listenBuildingUpgradeLoop(&b)
	b.activated_channel <- true

	for i := 0; i < 10; i++ {
		if !<-b.activated_send_channel {
			t.Error("Expected channel activated, was not")
		}
	}

	b.activated_send_done_channel <- true
}

func TestListenBuildingUpgrade(t *testing.T) {
	var b BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)

	go listenBuildingUpgrade(&b)

	select {
	case <-b.activated_send_channel:
		t.Error("Got unexpected building upgrade activated send signal")
	case <-time.After(CHANNEL_TIMEOUT):
	}

	b.activated_channel <- true

	for i := 0; i < 10; i++ {
		if !<-b.activated_send_channel {
			t.Error("Expected channel activated, was not")
		}
	}

	b.activated_send_done_channel <- true
	b.activated_done_channel <- true
}

func TestStartStopBuildingUpgrade(t *testing.T) {
	var b BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)
	StartBuildingUpgrade(&b)

	select {
	case <-b.activated_send_channel:
		t.Error("Got unexpected building upgrade activated send signal")
	case <-time.After(CHANNEL_TIMEOUT):
	}

	ActivateBuildingUpgrade(&b)

	for i := 0; i < 10; i++ {
		if !<-b.activated_send_channel {
			t.Error("Expected channel activated, was not")
		}
	}

	StopBuildingUpgrade(&b)
}

func TestGetAggregateUpgradeRatio(t *testing.T) {
	var b1 BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)
	var b2 BuildingUpgrade = MakeBuildingUpgrade(BUILDING_UPGRADE_NAME, BUILDING_UPGRADE_RATIO, BUILDING_UPGRADE_COST)

	var upgrades []*BuildingUpgrade = []*BuildingUpgrade{
		&b1,
		&b2,
	}

	var upgrade_ratio float64

	upgrade_ratio = GetAggregateUpgradeRatio(upgrades)
	if upgrade_ratio != 1 {
		t.Error(fmt.Sprintf("Expected upgrade ratio %e, got %e", 1, upgrade_ratio))
	}

	StartBuildingUpgrade(&b1)
	StartBuildingUpgrade(&b2)

	ActivateBuildingUpgrade(&b1)
	time.Sleep(EPOCH_TIME)

	upgrade_ratio = GetAggregateUpgradeRatio(upgrades)
	if upgrade_ratio != BUILDING_UPGRADE_RATIO {
		t.Error(fmt.Sprintf("Expected upgrade ratio %e, got %e", BUILDING_UPGRADE_RATIO, upgrade_ratio))
	}

	ActivateBuildingUpgrade(&b2)
	time.Sleep(EPOCH_TIME)

	upgrade_ratio = GetAggregateUpgradeRatio(upgrades)
	if upgrade_ratio != BUILDING_UPGRADE_RATIO*BUILDING_UPGRADE_RATIO {
		t.Error(fmt.Sprintf("Expected upgrade ratio %e, got %e", BUILDING_UPGRADE_RATIO*BUILDING_UPGRADE_RATIO, upgrade_ratio))
	}

	StopBuildingUpgrade(&b1)
	StopBuildingUpgrade(&b2)
}

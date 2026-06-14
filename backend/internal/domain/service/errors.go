package service

import "errors"

var (
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrInvalidToken              = errors.New("token has expired or has been already used")
	ErrUserNotFound              = errors.New("user not found")
	ErrIncorrectPassword         = errors.New("incorrect password")
	ErrInvalidPasswordFormat     = errors.New("invalid password format")
	ErrInvalidRefreshTokenFormat = errors.New("invalid refresh token format")
	ErrInvalidRefreshToken       = errors.New("token has expired or has been revoked")

	ErrInsufficientResources     = errors.New("insufficient resources")
	ErrBuildingLimitReached      = errors.New("building limit reached for this terrace level")
	ErrBuildingNotFound          = errors.New("building not found")
	ErrPositionOutOfBounds       = errors.New("position out of grid bounds")
	ErrPositionOccupied          = errors.New("position already occupied")
	ErrNotEnoughBuildingsInInventory = errors.New("not enough buildings in inventory")
	ErrBuildingNotPlaced         = errors.New("building not placed at this position")
	ErrUpgradeAlreadyActive      = errors.New("upgrade already in progress")
	ErrMaxLevelReached           = errors.New("building already at max level")

	ErrTroopNotFound             = errors.New("troop not found")
	ErrInsufficientArmyCapacity  = errors.New("insufficient army capacity")

	ErrCannotAttackSelf       = errors.New("cannot attack yourself")
	ErrDefenderNotFound       = errors.New("defender not found")
	ErrTerraceLevelMismatch   = errors.New("defender terrace level does not match")
	ErrDefenderShieldActive   = errors.New("defender has an active shield")
	ErrBattleNotFound         = errors.New("battle not found")
	ErrBattleNotPending       = errors.New("battle is not in pending state")
	ErrInvalidDeployment      = errors.New("invalid troop deployment")
	ErrInsufficientArmyTroops = errors.New("insufficient troops in army for deployment")
)

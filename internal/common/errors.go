package common

import "emperror.dev/errors"

const ServiceStartedNormallyButDoesNotBlock = errors.Sentinel("service has already been started")

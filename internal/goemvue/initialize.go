package goemvue

func (t *EmVueCloudSession) Initialize() error {
	/*
	 * Get a key
	 */
	err := t.authorize()
	if err != nil {
		return err
	}

	/*
	 * Get the basic customer information
	 */

	customer, err := t.getBasicCustomerInfo()
	if err != nil {
		return err
	}

	/*
	 * Load the device list
	 */

	fullCustomerInfo, err := t.getFullCustomerInfo(customer.CustomerGID)
	if err != nil {
		return err
	}
	
	t.CustomerInfo = fullCustomerInfo

	return nil
}

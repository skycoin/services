
#include "factory_test.h"
#include "trezor.h"
#include "oled.h"
#include "usb.h"
#include "buttons.h"
#include "layout.h"
#include "layout2.h"
#include "util.h"


void factoryTest(void)
{

	layoutDialog(&bmp_icon_question, ("Cancel"), ("Test Device"), NULL, ("Do you really want to"), ("test your Wallet?"), NULL, NULL, NULL, NULL);

	// wait until NoButton is released
	usbTiny(1);
	do {
		usbSleep(5);
		buttonUpdate();
	} while (!button.YesUp);

	// wait for confirmation/cancellation of the dialog
	do {
		usbSleep(5);
		buttonUpdate();
	} while (!button.YesUp && !button.NoUp);
	usbTiny(0);

	if (button.YesUp) { // if user say yes

		factoryTest_screen();

		if(factoryTest_buttons(0)&&factoryTest_buttons(1))
		{
			layoutDialogSwipe(&bmp_icon_info, NULL,"Exit",NULL,NULL, ("Buttons are OK"),NULL,NULL,NULL,NULL);
		}
		else
		{
			layoutDialogSwipe(&bmp_icon_info,NULL,"Exit",NULL,NULL, ("Buttons error"),NULL,NULL,NULL,NULL);
		}

		do { // Back home
			usbSleep(5);
			buttonUpdate();
		} while (!button.YesUp);
	}
	return;
}

void factoryTest_screen(void)
{
	oledClear();
	oledRefresh();

	for (int i = 2; i < 12; i++) {

		delay(0x3B9ACA00); // Wait 1 second

		if (i%2)
		{
		oledBox(0, 0, 127, 63, false);
		oledRefresh();
		}
		else
		{
		oledBox(0, 0, 127, 63, true);
		oledRefresh();
		}
	}
	return;
}

bool factoryTest_buttons(bool buttonToTest)
{
	bool testResult=0;

	layoutDialog(&bmp_icon_question,((buttonToTest) ? NULL:"No"),((buttonToTest) ? "Yes":NULL),NULL,((buttonToTest) ? "Press Yes":" Press No"),(""), NULL, NULL, NULL, NULL);

	do {
		usbSleep(5);
		buttonUpdate();
	} while (!button.YesUp && !button.NoUp);

	if ((button.NoUp && !buttonToTest) || (button.YesUp && buttonToTest) ) {
		// Button is ok
		testResult=1;
	} else{
		// User is stupid
		testResult=0;
	}

	oledClear();
	oledRefresh();

	return testResult;
}

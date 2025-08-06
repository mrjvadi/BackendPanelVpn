package event_handler

import (
	"fmt"
	"github.com/avct/uasurfer"
	"github.com/mrjvadi/BackendPanelVpn/events"
	type_event "github.com/mrjvadi/BackendPanelVpn/types/type-event"
	ptime "github.com/yaa110/go-persian-calendar"
	"gopkg.in/telebot.v4"
	"strings"
	"time"
)

// handleServiceJob processes job requests and simulates work.
func (b *BotEventHandler) eventLogin(evt events.Event) {
	if loginEvent, ok := evt.Payload.(type_event.EventBotLogin); ok {

		raw := strings.TrimSuffix(loginEvent.LoginTime, " UTC")

		// ۱. پارس کردن رشته‌ی زمان ورود به time.Time
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return
		}

		// ۲. تبدیل به زمان شمسی
		jTime := ptime.New(t)

		// ۳. فرمت کردن به شکل yyyy/MM/dd ساعت HH:mm:ss
		dateStr := fmt.Sprintf("%04d/%02d/%02d  %02d:%02d:%02d",
			jTime.Year(), jTime.Month(), jTime.Day(),
			jTime.Hour(), jTime.Minute(), jTime.Second(),
		)

		ua := uasurfer.Parse(loginEvent.DeviceInfo)

		// 3. استخراج اطلاعات مرورگر
		browserName := ua.Browser.Name.String() // مثلاً Chrome, Firefox, Safari, ...

		// 4. استخراج اطلاعات سیستم‌عامل
		osName := ua.OS.Name.String() // Windows, MacOS, Linux, Android, iOS, ...

		// 5. استخراج نوع دستگاه
		deviceType := ua.DeviceType.String() // Desktop, Phone, Tablet, Bot, Console, ...

		// ۴. ساخت پیام نهایی
		responseMessage := fmt.Sprintf(
			"🎉 سلام %s!\n\n"+
				"شما با موفقیت در تاریخ %s وارد شدید.\n"+
				"🖥️ مرورگر: %s\n"+
				"💻 سیستم‌عامل: %s\n"+
				"📱 نوع دستگاه: %s\n"+
				"🌐 آی‌پی: %s\n"+
				" شناسهٔ درخواست: %s\n"+
				"✅ کد وضعیت: %d",
			loginEvent.Username,
			dateStr,
			browserName,
			osName,
			deviceType,
			loginEvent.IPAddress,
			loginEvent.RequestID,
			loginEvent.Code,
		)

		// ۵. ارسال پیام
		b.bot.Send(telebot.ChatID(loginEvent.TelegramID), responseMessage)
	}
}

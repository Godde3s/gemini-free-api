<div dir="rtl">

# ✨ Gemini-Free-API

**پل رایگان و مستقیم به جمنای وب** (Gemini 3.6 Flash، 3.5 Flash-Lite و 3.1 Pro) — بدون کلید رسمی Google AI Studio، بدون مرورگر، فقط HTTP.
هر کلاینتی که با OpenAI یا Anthropic حرف می‌زند (Hermes، OpenCode، Cursor، Cline، Chatbox، LobeChat، …) مستقیم به این بریج وصل می‌شود.

> 🏗️ معماریِ برنده‌ی [GLM-Free-API](https://github.com/Godde3s/glm-free-api) و [Qwen-Free-API](https://github.com/Godde3s/qwen-free-api) — حالا برای gemini.google.com
>
> ⚡ برخلاف بقیه بریج‌های جمنای: **حالت مهمان واقعاً کار می‌کند** (زنده تست شده — حتی از IP دیتاسنتر)، مدل‌ها **زنده کشف** می‌شوند، استریم **افزایشی واقعی** است و ابزارها (tool calling) هم پشتیبانی می‌شوند.

---

## ✨ چرا این پروژه

| قابلیت | توضیح |
|---|---|
| 🚀 **تک‌فایل، تک‌دستور** | یک باینری، `./start.sh` یا `start.bat` — تمام. داشبورد فارسی داخل همان باینری embed شده |
| 👻 **حالت مهمان بدون کوکی** | بدون هیچ اکانتی کار می‌کند — `StreamGenerate` گوگل درخواست بی‌کوکی را با `at` خالی قبول می‌کند (زنده اثبات شد) |
| 🍪 **اکانت + چرخش خودکار کوکی** | کوکی‌های `__Secure-1PSID` را بدهید؛ درد اصلی بریج‌های جمنای یعنی انقضای مکرر `__Secure-1PSIDTS` به‌صورت خودکار حل می‌شود (RotateCookies + Set-Cookie + ذخیره روی دیسک) |
| 🔌 **دو پروتکل کامل** | `/v1/chat/completions` (OpenAI) + `/v1/messages` (Anthropic) با استریم SSE |
| 🧠 **تفکر زنده** | بلوک thoughts جمنای به `reasoning_content` استریم می‌شود (`reasoning:false` برای خاموش‌کردن) |
| 🛠️ **Tool Calling واقعی** | با `AGENT_MODE=true` ابزارهای OpenAI/Anthropic به prompt ترجمه و جواب‌ها به `tool_calls`/`tool_use` برمی‌گردند — برای OpenCode/Hermes/Cline |
| 🖼️ **ویژن واقعی** | `image_url` (data URI یا URL) → آپلود واقعی به `content-push.googleapis.com` → جمنای عکس را می‌بیند |
| 🔍 **کشف زنده مدل** | مدل‌ها با RPC داخلی گوگل (`otAQ7b`) در زمان اجرا کشف می‌شوند — مدل جدید گوگل بدون تغییر کد ظاهر می‌شود |
| 👥 **چند-اکانت واقعی** | round-robin + cooldown نمایی روی 429 + failover شفاف قبل از اولین بایت |
| 🪶 **سبک** | یک باینری ~10MB، بدون CDN، بدون دیتابیس، بدون Python، بدون curl-cffi |

## 📦 مدل‌ها

مدل‌ها **زنده** از اکانت شما کشف می‌شوند؛ این نام‌های پیش‌فرض وب جمنای هستند:

| مدل آیدی | مدل وب | کاربرد |
|---|---|---|
| **`gemini-3.6-flash`** ⭐ | 3.6 Flash | پیش‌فرض — سریع، برای اکثر کارها |
| `gemini-3.5-flash-lite` | 3.5 Flash-Lite | سبک‌ترین و سریع‌ترین |
| `gemini-3.1-pro` | 3.1 Pro | قوی‌ترین — کدنویسی و استدلال سنگین |

آلیاس‌ها (`gemini-flash`، `gemini-pro`، `flash`، `pro`، خود hex-id داخلی مثل `9d8ca3786ebdfbea`) همه پذیرفته می‌شوند. مدل ناشناخته → پیش‌فرض flash.

## 🚀 شروع در ۳۰ ثانیه

```bash
# 1) باینری را از Releases بگیرید یا با Go بسازید:
go build -trimpath -ldflags="-s -w" -o gemini-api .

# 2) اجرا:
./start.sh          # لینوکس/مک   — ویندوز: start.bat
```

همین! داشبورد روی `http://localhost:8080` باز می‌شود. بدون هیچ کوکی‌ای، حالت مهمان فعال است.

```bash
# 3) تست:
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer gemini" \
  -H "Content-Type: application/json" \
  -d '{"model":"gemini-3.6-flash","messages":[{"role":"user","content":"سلام! خودت رو معرفی کن"}]}'
```

## 🍪 اکانت گوگل (اختیاری ولی پیشنهادی)

حالت مهمان محدودیت سهمیه دارد و دسترسی مدل‌ها کمتر است. با یک کوکی ۶۰ ثانیه‌ای همه‌چیز باز می‌شود:

**راه آسان — ابزار خودکار:**
```bash
./gemini-login     # مرورگر باز می‌شود، لاگین کنید، کوکی‌ها خودکار ذخیره می‌شوند
```

**راه دستی:**
1. `gemini.google.com` → F12 → Application → Cookies
2. `__Secure-1PSID` و `__Secure-1PSIDTS` را کپی کنید
3. در `.env`:
```env
GEMINI_COOKIES=__Secure-1PSID=g.a00...; __Secure-1PSIDTS=ab12...
```

**چند-اکانت:**
```env
GEMINI_COOKIES=[{"psid":"...","psidts":"..."},{"psid":"...","psidts":"..."}]
```

> 💡 گوگل `__Secure-1PSIDTS` را مرتب می‌چرخاند و همین بزرگ‌ترین درد بریج‌های جمنای است. این بریج مقدار تازه را از پاسخ‌ها می‌گیرد، در `gemini-cookies-cache.json` ذخیره می‌کند و بعد از ری‌استارت هم استفاده می‌کند؛ در خطای احراز هویت یک‌بار `RotateCookies` گوگل را امتحان می‌کند و فقط بعد از آن اکانت را مرده اعلام می‌کند.

## 🤖 اتصال ایجنت‌ها (OpenCode / Hermes / Cline)

**OpenAI-compatible:**
```env
OPENAI_BASE_URL=http://localhost:8080/v1
OPENAI_API_KEY=gemini
MODEL=gemini-3.6-flash
```

**Anthropic-compatible:**
```env
ANTHROPIC_BASE_URL=http://localhost:8080
ANTHROPIC_API_KEY=gemini
```

برای ایجنت‌های ابزار-محور حتماً:
```env
AGENT_MODE=true
```
ابزارهای OpenAI/Anthropic به یک پرامپت ساخت‌یافته ترجمه می‌شوند، مدل بلوک `<<<TOOL_CALL>>>` می‌سازد و بریج همان را به `tool_calls`/`tool_use` استاندارد برمی‌گرداند — گردش کار R1→R2 کامل تست شده است.

## 🖼️ ویژن (عکس)

```json
{
  "model": "gemini-3.6-flash",
  "messages": [{
    "role": "user",
    "content": [
      {"type": "text", "text": "این عکس چیه؟"},
      {"type": "image_url", "image_url": {"url": "https://.../photo.jpg"}}
    ]
  }]
}
```
عکس واقعاً به سرور گوگل آپلود می‌شود (endpoint رسمی وب اپ) و جمنای آن را می‌بیند — نه متن جایگزین.

## 🎛️ داشبورد

روی `/`: وضعیت زنده اکانت‌ها (سالم/استراحت با شمارش معکوس/مرده)، کپی اندپوینت‌ها، اسنیپت آماده curl/Python/Node/ایجنت، و **پلی‌گراند استریم** برای تست زنده با یک کلیک.

## 🐳 Docker

```bash
docker compose up -d
# یا
docker build -t gemini-free-api . && docker run -d -p 8080:8080 -v $(pwd)/.env:/app/.env gemini-free-api
```

## ⚙️ همه متغیرهای محیطی

| متغیر | پیش‌فرض | توضیح |
|---|---|---|
| `GEMINI_COOKIES` | خالی | کوکی اکانت‌ها (سه فرمت — بالا را ببینید). خالی = حالت مهمان |
| `GEMINI_COOKIES_FILE` | `gemini-cookies.json` | فایل کوکی‌ها (خروجی `gemini-login`) |
| `PORT` / `HOST` | `8080` / `0.0.0.0` | سرور |
| `AUTH_TOKEN` | `gemini` | کلید کلاینت‌ها |
| `THINKING` | `true` | پیش‌فرض استریم reasoning |
| `AGENT_MODE` | `false` | ترجمه ابزار برای ایجنت‌ها |
| `ACCOUNT_COOLDOWN_BASE` | `120` | اولین cooldown روی 429 (ثانیه، نمایی ×۲) |
| `ACCOUNT_QUEUE_TIMEOUT` | `120` | حداکثر انتظار صف وقتی همه اکانت‌ها لیمیت‌اند |
| `INIT_TOKEN_TTL` | `1800` | کش SNlM0e/bl/f.sid هر اکانت |

## 🧯 عیب‌یابی

| خطا | معنی و راه‌حل |
|---|---|
| «کوکی‌های اکانت گوگل منقضی شده‌اند» | دوباره `gemini-login` اجرا کنید یا کوکی‌ها را تازه کنید |
| «سقف درخواست گوگل (HTTP 429)» | سهمیه پر شده — چند دقیقه صبر یا اکانت بیشتر |
| «IP شما موقتاً توسط گوگل فلگ شده (کد ۱۰۶۰)» | IP دیتاسنتر/VPN فلگ شده — پروکسی عوض کنید (`HTTPS_PROXY` پشتیبانی می‌شود) |
| «سقف استفاده این مدل (کد ۱۰۳۷)» | سهمیه همان مدل پر شده — مدل دیگری بزنید |
| پاسخ نمی‌آید | `LOG_LEVEL=debug` و لاگ را ببینید |

## 📁 ساختار پروژه

```
├── main.go                  # نقطه ورود
├── internal/gbridge/
│   ├── gemini.go            # کلاینت آپستریم (init/StreamGenerate/دلتا/خطاها)
│   ├── frames.go            # پارسر فریم‌های length-prefixed گوگل (UTF-16-aware)
│   ├── accounts.go          # پول کوکی + چرخش PSIDTS + cooldown
│   ├── models.go            # کشف زنده مدل (batchexecute otAQ7b)
│   ├── images.go            # ویژن واقعی (content-push upload)
│   ├── agent.go / agent_legacy.go  # shim ابزارها
│   ├── anthropic.go         # /v1/messages
│   ├── handlers.go          # /v1/chat/completions
│   └── dashboard.go         # داشبورد فارسی (embed)
└── cmd/gemini-login/        # ابزار capture کوکی با playwright
```

## ⚠️ نکات مهم

- این پروژه یک بریج **شخصی** به وب‌اپ جمنای است؛ ToS گوگل را مرور کنید و در حجم بالا استفاده نکنید.
- کوکی‌ها فقط روی دیسک خودتان ذخیره می‌شوند و به هیچ سروری فرستاده نمی‌شوند.
- اگر سرور را عمومی می‌کنید حتماً `AUTH_TOKEN` را عوض کنید.
- **کوکی‌ها را با کسی به اشتراک نگذارید** — دسترسی کامل اکانت گوگل می‌دهد.

## 🇬🇧 English Summary

Free OpenAI & Anthropic-compatible bridge for the Gemini **web app** (gemini.google.com). Single static Go binary with an embedded Persian dashboard. Works **cookie-less in guest mode** (verified live, even from datacenter IPs) or with `__Secure-1PSID`/`__Secure-1PSIDTS` account cookies featuring **automatic PSIDTS rotation** (the #1 pain point of Gemini bridges). Live model discovery via Google's internal RPC, real incremental SSE streaming from Google's length-prefixed frames, real vision through content-push uploads, optional tool-calling shim for agent frameworks, multi-account pool with exponential 429 cooldown and pre-byte failover. Run: `./start.sh` (or `start.bat` on Windows) → dashboard at `http://localhost:8080`. See `.env.example` for every option.

</div>

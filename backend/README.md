> **إلزامي عند العمل على ميزة في الباك إند:** اقرأ ملفها في [`../docs/features/`](../docs/features/README.md) قبل التنفيذ، وحدّثه في نفس المهمة بعد أي تغيير. راجع [`../AGENTS.md`](../AGENTS.md).

# واجهة قارئ البرمجية

خادم REST بلغة Go يستخدم MySQL 8 مع InnoDB و`utf8mb4`. لا يحتاج الفرونت إلى اتصال مباشر بقاعدة البيانات. تبقى نقاط REST الـ11 كما كانت.

## التشغيل المحلي

على جهاز Windows الحالي، خدمة `MySQL80` تستخدم المنفذ `3307` بينما `3306` مخصص لـMariaDB. يمكن تشغيل `scripts/setup-local-mysql.ps1` مرة واحدة من PowerShell بصلاحية مسؤول لإعادة تعيين كلمة مرور MySQL المحلية، وإنشاء قاعدة وحساب منفصلين للمنصة، وكتابة `backend/.env`؛ سيطلب كلمة المرور الجديدة داخل النافذة ولا يطبعها. تبقى التعليمات اليدوية أدناه صالحة لخادم MySQL آخر، مع تعديل المنفذ حسب بيئته.

1. أنشئ قاعدة فارغة ومستخدمًا مخصصًا لها في MySQL:

```sql
CREATE DATABASE qare CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'qare'@'127.0.0.1' IDENTIFIED BY 'change-me';
GRANT ALL PRIVILEGES ON qare.* TO 'qare'@'127.0.0.1';
```

2. من `backend/`، اضبط متغيرات البيئة وشغّل الخادم:

```powershell
$env:MYSQL_DSN='qare:change-me@tcp(127.0.0.1:3306)/qare'
$env:FRONTEND_ORIGIN='http://localhost:3000'
$env:PORT='8080'
go mod download
go run .
```

ملف `.env.example` مرجع للقيم؛ Go لا يحمّل `.env` تلقائيًا. إذا احتوت كلمة المرور على رموز خاصة، استخدم صيغة DSN المتوافقة مع `go-sql-driver/mysql`. في النشر اضبط `MYSQL_DSN` و`FRONTEND_ORIGIN` من أسرار وبيئة الاستضافة؛ استخدم اتصال TLS بقاعدة البيانات إذا كان الاتصال عبر شبكة خارجية. أنشئ قاعدة البيانات أولًا؛ الخادم ينشئ الجداول ويضع ستة كتب ومراجعتين فقط عندما يكون جدول الكتب فارغًا. حساب المراجعات التجريبية للعرض ولا يمكن الدخول به.

اضبط `NEXT_PUBLIC_API_URL` في `frontend/` إلى عنوان الباك إند العام مع `/api`، مثل `https://api.example.com/api`. يجب أن يطابق `FRONTEND_ORIGIN` أصل موقع الفرونت الفعلي، مثل `https://example.com`.

## النشر على Render

لربط خدمة Aiven، نزّل شهادة CA من صفحة الخدمة، ثم شغّل `backend/scripts/prepare-aiven.ps1` من جذر المشروع في PowerShell. يجمع المساعد بيانات الاتصال محليًا، ويحفظ كلمة المرور مشفرة لحساب Windows الحالي في `.cache/aiven-connection.json`، ويضع الشهادة في `backend/aiven-ca.pem`. كلا الملفين مستبعد من Git. تُضبط `MYSQL_DSN` و`MYSQL_CA_CERT_PATH` وشهادة `mysql-ca.pem` في Render بعد الموافقة على إرسال بيانات الاتصال للخدمة.

ملف [`../render.yaml`](../render.yaml) يجهز خدمة Go من مجلد `backend/` مع فحص صحة `GET /api/books`. اربط مستودع GitHub بميزة Blueprint في Render، وأدخل `MYSQL_DSN` لقاعدة **MySQL مستضافة** و`FRONTEND_ORIGIN` لرابط Vercel النهائي. لا تستخدم عنوان MySQL المحلي `127.0.0.1:3307` داخل Render؛ سيشير إلى حاوية Render نفسها. إذا لم تُنشأ قاعدة مستضافة بعد، أعدّها وانقل بياناتك إليها قبل اعتبار النشر مكتملًا. الخادم يأخذ `PORT` من Render تلقائيًا.

## نقل بيانات SQLite القديمة

احتفظ بنسخة احتياطية من `qare.db`، وأوقف خادم Go القديم قبل النقل حتى لا تتغير البيانات أثناء النسخ. أنشئ قاعدة MySQL **فارغة** ولا تشغّل API الجديد عليها قبل النقل، ثم من `backend/`:

```powershell
$env:MYSQL_DSN='qare:change-me@tcp(127.0.0.1:3306)/qare'
go run ./cmd/importsqlite ./qare.db
```

الأداة تنشئ المخطط، ثم تنقل `users`, `books`, `sessions`, `reviews`, `library` داخل معاملة واحدة مع الحفاظ على المعرفات وتجزئات كلمات المرور والرموز. ترفض النقل إذا وُجدت بيانات في MySQL، حتى لا تكتب فوق بيانات موجودة. بعد نجاحها شغّل `go run .`؛ لن تُضاف بيانات أولية فوق الكتب المنقولة. لا تحذف ملف SQLite قبل التحقق من الحسابات والمراجعات والمكتبة في النسخة الجديدة.

## تنظيم الكود

| الموضع | المسؤولية |
|---|---|
| `main.go` | قراءة البيئة وفتح MySQL وتشغيل HTTP |
| `internal/routes/routes.go` | المسارات وCORS |
| `internal/handlers/` | طلبات HTTP والاستجابات |
| `internal/services/` | التحقق وقواعد العمل |
| `internal/repositories/` | SQL والترحيلات والبيانات الأولية |
| `cmd/importsqlite/` | أداة نقل البيانات القديمة فقط؛ هي المستعمل الوحيد لمكتبة SQLite |

الطلب يمر من route إلى handler، ثم service عند وجود قواعد عمل، ثم repository إلى MySQL. القراءات البسيطة تستدعي repository مباشرة من handler. أضف ترحيلًا جديدًا إلى `internal/repositories/migrations.go` عند تغيير المخطط ولا تعدّل ترحيلًا طُبّق. يستخدم `schema_migrations` لتتبع الإصدار، وMySQL لا يجعل DDL جزءًا من معاملة قابلة للتراجع؛ لذا يجب أن تكون عبارات الترحيل قابلة لإعادة التشغيل. يمنع قفل MySQL المسمى تنفيذ الترحيلات بالتزامن عند تشغيل عدة نسخ.

## الاختبارات

```powershell
go test ./...
go vet ./...
```

يتخطى الاختبار التكاملي عند غياب `MYSQL_TEST_DSN`. لاختبار REST والترحيلات فعليًا، أنشئ قاعدة MySQL **منفصلة مخصصة للاختبار** واضبط `MYSQL_TEST_DSN` ثم أعد `go test ./...`. اختبار REST ينشئ حسابًا مؤقتًا ويحذفه بعده.

## REST API

جميع المسارات تبدأ بـ `/api`. أرسل `Content-Type: application/json`. المسارات المحمية تحتاج `Authorization: Bearer <token>`.

| الطريقة | المسار | الوظيفة / الجسم |
|---|---|---|
| POST | `/auth/register` | `{name,email,password}`؛ ينشئ مستخدمًا ويعيد `{user,token}` |
| POST | `/auth/login` | `{email,password}`؛ يعيد `{user,token}` |
| GET | `/auth/me` | بيانات المستخدم الحالي |
| POST | `/auth/logout` | إبطال الجلسة |
| GET | `/books?q=` | قائمة الكتب؛ البحث بالعنوان أو المؤلف |
| GET | `/books/{id}` | تفاصيل الكتاب ومراجعاته |
| PUT | `/books/{id}/review` | إنشاء/تعديل مراجعتك: `{rating: 1..5, body}` |
| DELETE | `/books/{id}/review` | حذف مراجعتك |
| GET | `/library` | كتب المستخدم بتصنيفاتها |
| PUT | `/library/{id}` | إضافة/نقل كتاب: `{status: "want"|"reading"|"read"}` |
| DELETE | `/library/{id}` | إزالة كتاب من المكتبة |

أخطاء الطلب تعيد JSON بصيغة `{ "error": "..." }`.

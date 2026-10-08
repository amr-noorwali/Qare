> **إلزامي عند العمل على ميزة في الباك إند:** اقرأ ملفها في [`../docs/features/`](../docs/features/README.md) قبل التنفيذ، وحدّثه في نفس المهمة بعد أي تغيير. راجع [`../AGENTS.md`](../AGENTS.md).

# واجهة قارئ البرمجية

خادم REST بلغة Go مع SQLite. قاعدة البيانات المحلية هي `qare.db` داخل هذا المجلد عند التشغيل منه. الملف موجود ومجهز للتجربة، ويُنشئ الخادم الجداول والكتب التجريبية ومراجعتين تلقائيًا إذا لم يكن موجودًا. حساب المراجعات التجريبية للعرض فقط؛ أنشئ حسابًا لتجربة الكتابة.

## التشغيل

```powershell
Copy-Item .env.example .env
go mod download
go run .
```

يقرأ الخادم متغيرات البيئة من بيئة التشغيل؛ ملف `.env` مرجع فقط، ولا يُحمّل تلقائيًا. في PowerShell: `$env:PORT='8080'; $env:DATABASE_PATH='./qare.db'; $env:FRONTEND_ORIGIN='http://localhost:3000'`. القيم المذكورة افتراضية. استخدم `go test ./...` للتحقق.

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

أخطاء الطلب تعيد JSON بصيغة `{ "error": "..." }`. قاعدة البيانات محفوظة في `DATABASE_PATH`؛ احذف الملف لإعادة البيانات الأولية.

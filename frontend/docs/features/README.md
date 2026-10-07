> **إلزامي:** قبل تعديل ميزة في الفرونت، اقرأ ملفها هنا والتوثيق الشامل المناظر في `../../../docs/features/`. بعد تنفيذ المهمة، حدّث الاثنين في نفس المهمة. إذا مست المهمة ميزتين، حدّث ملفيهما.

# توثيق ميزات الفرونت إند

هذه الملفات داخل مشروع الواجهة نفسه. تصف ما تنفذه صفحات Next.js والمكونات وخدمة REST، بينما [التوثيق الشامل](../../../docs/features/README.md) يشرح أيضًا Go وSQLite.

| الميزة | ملف الفرونت | الصفحة أو موضعها |
|---|---|---|
| الحساب والجلسة | [authentication.md](authentication.md) | `/login`, `/register`، الترويسة |
| الكتب والبحث | [catalog-search.md](catalog-search.md) | `/` |
| تفاصيل الكتاب | [book-details.md](book-details.md) | `/books/[id]` |
| المراجعات والتقييم | [reviews-ratings.md](reviews-ratings.md) | داخل `/books/[id]` |
| المكتبة الشخصية | [personal-library.md](personal-library.md) | `/library` وجزء من `/books/[id]` |
| الهوية البصرية | [visual-identity.md](visual-identity.md) | مشتركة بين كل الصفحات |

لفهم توزيع الملفات بين `app/` و`components/` و`lib/` راجع [خريطة الفرونت](../../../docs/features/frontend-map.md). ليس لدينا مجلد كود `features/` مستقل؛ هذا المجلد للتوثيق.

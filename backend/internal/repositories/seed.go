package repositories

import "database/sql"

func seed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM books").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	books := [][5]string{
		{"موسم الهجرة إلى الشمال", "الطيب صالح", "رواية عن العودة والهوية واللقاء المعقد بين عالمين، تُروى بصوت شاب يعود إلى قريته على ضفاف النيل.", "https://covers.openlibrary.org/b/isbn/9780141187198-L.jpg", "رواية"},
		{"رجال في الشمس", "غسان كنفاني", "ثلاث حكايات تلتقي في رحلة بحث عن حياة أفضل، في رواية عربية قصيرة ومؤثرة.", "https://covers.openlibrary.org/b/isbn/9780894103285-L.jpg", "أدب عربي"},
		{"الخيميائي", "باولو كويلو", "رحلة راعٍ شاب وراء حلمه، يقوده الطريق إلى أسئلة عن القدر والشجاعة والإصغاء إلى القلب.", "https://covers.openlibrary.org/b/isbn/9780062315007-L.jpg", "رواية"},
		{"عالم صوفي", "جوستاين غاردر", "مدخل روائي إلى تاريخ الفلسفة، يبدأ برسائل غامضة تصل إلى فتاة صغيرة.", "https://covers.openlibrary.org/b/isbn/9780374530716-L.jpg", "فلسفة"},
		{"الأمير الصغير", "أنطوان دو سانت إكزوبيري", "حكاية شاعرية عن الصداقة والخيال وما يغفله الكبار في زحام الحياة.", "https://covers.openlibrary.org/b/isbn/9780156012195-L.jpg", "كلاسيكيات"},
		{"1984", "جورج أورويل", "رواية عن السلطة والمراقبة واللغة في عالم مستقبلي شديد القسوة.", "https://covers.openlibrary.org/b/isbn/9780451524935-L.jpg", "خيال ديستوبي"},
	}
	for _, book := range books {
		if _, err := tx.Exec("INSERT INTO books(title,author,description,cover_url,genre) VALUES(?,?,?,?,?)", book[0], book[1], book[2], book[3], book[4]); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`INSERT INTO users(name,email,password_hash) VALUES('قارئ من منصة قارئ','demo@qare.local','seed-only');
 INSERT INTO reviews(book_id,user_id,rating,body) VALUES(1,1,5,'لغة آسرة وأسئلة تبقى معك بعد الصفحة الأخيرة.'),(3,1,4,'رحلة خفيفة وملهمة للعودة إلى الأحلام.');`); err != nil {
		return err
	}
	return tx.Commit()
}

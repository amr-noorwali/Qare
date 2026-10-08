package repositories

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`PRAGMA foreign_keys=ON;
 CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS sessions (token TEXT PRIMARY KEY, user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, expires_at TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, title TEXT NOT NULL, author TEXT NOT NULL, description TEXT NOT NULL, cover_url TEXT NOT NULL, genre TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS reviews (id INTEGER PRIMARY KEY, book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE, user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, rating INTEGER NOT NULL CHECK(rating BETWEEN 1 AND 5), body TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UNIQUE(book_id,user_id));
 CREATE TABLE IF NOT EXISTS library (user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE, status TEXT NOT NULL CHECK(status IN ('want','reading','read')), PRIMARY KEY(user_id,book_id));
 CREATE INDEX IF NOT EXISTS idx_reviews_book_created ON reviews(book_id,created_at DESC);
 CREATE INDEX IF NOT EXISTS idx_library_book ON library(book_id);
 CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);`)
	if err != nil {
		db.Close()
		return nil, err
	}
	if err = seed(db); err != nil {
		db.Close()
		return nil, err
	}
	if _, err = db.Exec("UPDATE users SET name=? WHERE email=? AND password_hash=?", "قارئ من منصة قارئ", "demo@qare.local", "seed-only"); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func seed(db *sql.DB) error {
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
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
	for _, b := range books {
		if _, err := db.Exec("INSERT INTO books(title,author,description,cover_url,genre) VALUES(?,?,?,?,?)", b[0], b[1], b[2], b[3], b[4]); err != nil {
			return err
		}
	}
	_, err := db.Exec(`INSERT INTO users(name,email,password_hash) VALUES('قارئ من منصة قارئ','demo@qare.local','seed-only');
 INSERT INTO reviews(book_id,user_id,rating,body) VALUES(1,1,5,'لغة آسرة وأسئلة تبقى معك بعد الصفحة الأخيرة.'),(3,1,4,'رحلة خفيفة وملهمة للعودة إلى الأحلام.');`)
	return err
}

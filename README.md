# WhatsApp AI Call Assistant Bot (Golang)

Aplikasi bot panggilan suara WhatsApp otomatis berbasis Go (Golang) yang ditenagai oleh:
- **`whatsmeow` & `meowcaller`**: Protocol client WhatsApp VOIP call murni di Go (mengangkat & menelpon, enkoding/dekoding audio MLow codec).
- **Groq API (`whisper-large-v3`)**: Transkripsi suara (Speech-to-Text) ultra cepat.
- **Google Gemini API (`gemini-3.6-flash`)**: AI percakapan generasi terbaru (Gemini 3.6) yang cepat dan hemat token untuk panggilan suara.
- **Anime TTS & Google Translate TTS**: Sintesis suara (Text-to-Speech) yang diputar langsung di dalam saluran telepon WhatsApp.

---

## 📌 Fitur Utama (Features)

1. **Model AI Terbaru (`gemini-3.6-flash`)**: Menggunakan Gemini 3.6 dengan dukungan fallback otomatis ke model Gemini versi lain jika diperlukan.
2. **Fitur Pengaman / Fitur Owner (`OWNER`)**: Membatasi akses bot telepon dan perintah chat hanya untuk nomor HP yang terdaftar di `.env`.
3. **Dukungan Perintah di Grup & DM**: Perintah `!aicall 628xxx` / `.call 628xxx` mendukung berbagai jenis pesan WhatsApp (Group Message, Ephemeral, ViewOnce, Extended Text).
4. **Validasi Input Nomor**: Menampilkan peringatan penggunaan jika perintah dipanggil tanpa menyebutkan nomor target.
5. **System Prompt, TTS Speed & Language via `.env`**: Mengatur gaya bicara AI, kecepatan suara (`TTS_SPEED`), bahasa (`TTS_LANG`), dan karakter suara (`TTS_CHARACTER`) secara fleksibel.

---

## 📋 Prasyarat (Prerequisites)

- **Go (Golang)**: Versi 1.22 atau lebih baru.
- **API Keys (Gratis)**:
  - **Gemini API Key**: Dapatkan gratis di [Google AI Studio](https://aistudio.google.com/app/apikey)
  - **Groq API Key**: Dapatkan gratis di [Groq Console](https://console.groq.com/keys)

---

## 🚀 Cara Install & Menjalankan (Installation & Usage)

### 1. Clone & Masuk ke Direktori
```bash
git clone https://github.com/krsna081/assisten-ai-call.git
cd assisten-ai-call
```

### 2. Konfigurasi Environment Variable (`.env`)
Salin file template `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

Buka file `.env` dan atur konfigurasi Anda:
```env
GEMINI_API=AIzaSy...               # API Key Gemini Anda
GEMINI_MODEL=gemini-3.6-flash       # Model Gemini (default: gemini-3.6-flash)
GROQ_API=gsk_...                    # API Key Groq Anda

# Nomor Owner Bot (Dipisahkan koma jika lebih dari satu. Kosongkan jika ingin publik)
OWNER=628123456789,628987654321,604312345678

# System Prompt AI (Bisa disesuaikan sesuka Anda)
SYSTEM_PROMPT="Kamu adalah asisten suara AI pintar..."

# Karakter Suara, Bahasa & Kecepatan TTS
TTS_CHARACTER=特别周 Special Week (Umamusume Pretty Derby)
TTS_LANG=Mix
TTS_SPEED=1.0
```

### 3. Build dan Jalankan Aplikasi
```bash
go build -o ai-call .
./ai-call
```

---

## 📱 Petunjuk Penggunaan

1. **Perintah Chat (`!aicall` / `.call`)**
   - Di DM atau Grup WhatsApp, ketik:
     ```text
     !aicall 628123456789
     ```
   - Jika nomor target tidak diisi, bot akan merespon dengan instruksi format yang benar:
     `Gunakan format: !aicall 628xxxxxxxx`

2. **Pengamanan Owner (`OWNER`)**
   - Jika variabel `OWNER` di `.env` diisi, pengguna di luar daftar tersebut yang menelpon bot akan ditolak otomatis (*rejected*), dan penggunaan perintah `!aicall` akan ditolak.

---

## 🌐 Daftar Bahasa TTS (`TTS_LANG`)

Anda dapat memilih bahasa pengucapan untuk sintetis suara Anime TTS di `.env`:
- `Mix` *(Default: Campuran Bahasa)*
- `English` *(Bahasa Inggris)*
- `日本語` *(Bahasa Jepang)*
- `简体中文` *(Bahasa Mandarin / Chinese)*

---

## 🗣️ Daftar Lengkap Karakter Text-to-Audio (`TTS_CHARACTER`)

Salin teks karakter di bawah ini dan tempelkan pada variabel `TTS_CHARACTER` di file `.env`:

### 🐎 Umamusume Pretty Derby
- `特别周 Special Week (Umamusume Pretty Derby)` *(Default)*
- `无声铃鹿 Silence Suzuka (Umamusume Pretty Derby)`
- `东海帝王 Tokai Teio (Umamusume Pretty Derby)`
- `丸善斯基 Maruzensky (Umamusume Pretty Derby)`
- `富士奇迹 Fuji Kiseki (Umamusume Pretty Derby)`
- `小栗帽 Oguri Cap (Umamusume Pretty Derby)`
- `黄金船 Gold Ship (Umamusume Pretty Derby)`
- `伏特加 Vodka (Umamusume Pretty Derby)`
- `大和赤骥 Daiwa Scarlet (Umamusume Pretty Derby)`
- `大树快车 Taiki Shuttle (Umamusume Pretty Derby)`
- `草上飞 Grass Wonder (Umamusume Pretty Derby)`
- `菱亚马逊 Hishi Amazon (Umamusume Pretty Derby)`
- `目白麦昆 Mejiro Mcqueen (Umamusume Pretty Derby)`
- `神鹰 El Condor Pasa (Umamusume Pretty Derby)`
- `好歌剧 T.M. Opera O (Umamusume Pretty Derby)`
- `成田白仁 Narita Brian (Umamusume Pretty Derby)`
- `鲁道夫象征 Symboli Rudolf (Umamusume Pretty Derby)`
- `气槽 Air Groove (Umamusume Pretty Derby)`
- `爱丽数码 Agnes Digital (Umamusume Pretty Derby)`
- `青云天空 Seiun Sky (Umamusume Pretty Derby)`
- `玉藻十字 Tamamo Cross (Umamusume Pretty Derby)`
- `美妙 姿势 Fine Motion (Umamusume Pretty Derby)`
- `琵琶晨光 Biwa Hayahide (Umamusume Pretty Derby)`
- `重炮 Mayano Topgun (Umamusume Pretty Derby)`
- `曼城茶座 Manhattan Cafe (Umamusume Pretty Derby)`
- `美普波旁 Mihono Bourbon (Umamusume Pretty Derby)`
- `目白雷恩 Mejiro Ryan (Umamusume Pretty Derby)`
- `雪之美人 Yukino Bijin (Umamusume Pretty Derby)`
- `米浴 Rice Shower (Umamusume Pretty Derby)`
- `艾尼斯风神 Ines Fujin (Umamusume Pretty Derby)`
- `爱丽速子 Agnes Tachyon (Umamusume Pretty Derby)`
- `爱慕织姬 Admire Vega (Umamusume Pretty Derby)`
- `稻荷一 Inari One (Umamusume Pretty Derby)`
- `胜利奖券 Winning Ticket (Umamusume Pretty Derby)`
- `空中神宫 Air Shakur (Umamusume Pretty Derby)`
- `荣进闪耀 Eishin Flash (Umamusume Pretty Derby)`
- `真机伶 Curren Chan (Umamusume Pretty Derby)`
- `川上公主 Kawakami Princess (Umamusume Pretty Derby)`
- `黄金城市 Gold City (Umamusume Pretty Derby)`
- `樱花进王 Sakura Bakushin O (Umamusume Pretty Derby)`
- `采珠 Seeking the Pearl (Umamusume Pretty Derby)`
- `新光风 Shinko Windy (Umamusume Pretty Derby)`
- `东商变革 Sweep Tosho (Umamusume Pretty Derby)`
- `超级小溪 Super Creek (Umamusume Pretty Derby)`
- `醒目飞鹰 Smart Falcon (Umamusume Pretty Derby)`
- `荒 漠英雄 Zenno Rob Roy (Umamusume Pretty Derby)`
- `东瀛佐敦 Tosen Jordan (Umamusume Pretty Derby)`
- `中山庆典 Nakayama Festa (Umamusume Pretty Derby)`
- `成田大进 Narita Taishin (Umamusume Pretty Derby)`
- `西野花 Nishino Flower (Umamusume Pretty Derby)`
- `春乌拉拉 Haru Urara (Umamusume Pretty Derby)`
- `青竹回忆 Bamboo Memory (Umamusume Pretty Derby)`
- `待兼福来 Matikane Fukukitaru (Umamusume Pretty Derby)`
- `名将怒涛 Meisho Doto (Umamusume Pretty Derby)`
- `目白多伯 Mejiro Dober (Umamusume Pretty Derby)`
- `优秀素质 Nice Nature (Umamusume Pretty Derby)`
- `帝 王光环 King Halo (Umamusume Pretty Derby)`
- `待兼诗歌剧 Matikane Tannhauser (Umamusume Pretty Derby)`
- `生野狄杜斯 Ikuno Dictus (Umamusume Pretty Derby)`
- `目白善信 Mejiro Palmer (Umamusume Pretty Derby)`
- `大 拓太阳神 Daitaku Helios (Umamusume Pretty Derby)`
- `双涡轮 Twin Turbo (Umamusume Pretty Derby)`
- `里见光钻 Satono Diamond (Umamusume Pretty Derby)`
- `北部玄驹 Kitasan Black (Umamusume Pretty Derby)`
- `樱花千代 王 Sakura Chiyono O (Umamusume Pretty Derby)`
- `天狼星象征 Sirius Symboli (Umamusume Pretty Derby)`
- `目白阿尔丹 Mejiro Ardan (Umamusume Pretty Derby)`
- `八重无敌 Yaeno Muteki (Umamusume Pretty Derby)`
- `鹤丸刚志 Tsurumaru Tsuyoshi (Umamusume Pretty Derby)`
- `目白光明 Mejiro Bright (Umamusume Pretty Derby)`
- `樱花桂冠 Sakura Laurel (Umamusume Pretty Derby)`
- `成田路 Narita Top Road (Umamusume Pretty Derby)`
- `也文摄 辉 Yamanin Zephyr (Umamusume Pretty Derby)`
- `真弓快车 Aston Machan (Umamusume Pretty Derby)`
- `骏川手纲 Hayakawa Tazuna (Umamusume Pretty Derby)`
- `小林历奇 Kopano Rickey (Umamusume Pretty Derby)`
- `奇锐骏 Wonder Acute (Umamusume Pretty Derby)`
- `秋川理事长 President Akikawa (Umamusume Pretty Derby)`

### ⚔️ Genshin Impact
- `派蒙 Paimon (Genshin Impact)`
- `雷电将军 Raiden Shogun (Genshin Impact)`
- `胡桃 Hu Tao (Genshin Impact)`
- `甘雨 Ganyu (Genshin Impact)`
- `纳西妲 Nahida (Genshin Impact)`
- `刻晴 Keqing (Genshin Impact)`
- `神里绫华 Kamisato Ayaka (Genshin Impact)`
- `钟离 Zhongli (Genshin Impact)`
- `温迪 Venti (Genshin Impact)`
- `可莉 Klee (Genshin Impact)`
- `九条裟罗 Kujou Sara (Genshin Impact)`
- `芭芭 拉 Barbara (Genshin Impact)`
- `荒泷一 斗 Arataki Itto (Genshin Impact)`
- `早柚 Sayu (Genshin Impact)`
- `香菱 Xiangling (Genshin Impact)`
- `重云 Chongyun (Genshin Impact)`
- `流浪者 Wanderer (Genshin Impact)`
- `优菈 Eula (Genshin Impact)`
- `凝光 Ningguang (Genshin Impact)`
- `枫原万叶 Kaedehara Kazuha (Genshin Impact)`
- `赛诺 Cyno (Genshin Impact)`
- `诺艾尔 Noelle (Genshin Impact)`
- `八重神子 Yae Miko (Genshin Impact)`
- `凯亚 Kaeya (Genshin Impact)`
- `魈 Xiao (Genshin Impact)`
- `托马 Thoma (Genshin Impact)`
- `迪卢克 Diluc (Genshin Impact)`
- `夜兰 Yelan (Genshin Impact)`
- `鹿野院平藏 Shikanoin Heizou (Genshin Impact)`
- `辛焱 Xinyan (Genshin Impact)`
- `丽莎 Lisa (Genshin Impact)`
- `云堇 Yun Jin (Genshin Impact)`
- `坎蒂丝 Candace (Genshin Impact)`
- `罗莎莉亚 Rosaria (Genshin Impact)`
- `北斗 Beidou (Genshin Impact)`
- `珊瑚宫心海 Sangonomiya Kokomi (Genshin Impact)`
- `烟绯 Yanfei (Genshin Impact)`
- `久岐忍 Kuki Shinobu (Genshin Impact)`
- `宵宫 Yoimiya (Genshin Impact)`
- `安柏 Amber (Genshin Impact)`
- `迪奥娜 Diona (Genshin Impact)`
- `班尼特 Bennett (Genshin Impact)`
- ` 雷泽 Razor (Genshin Impact)`
- `阿贝多 Albedo (Genshin Impact)`
- `空 Player Male (Genshin Impact)`
- `神里绫人 Kamisato Ayato (Genshin Impact)`
- `琴 Jean (Genshin Impact)`
- `艾尔海 森 Alhaitham (Genshin Impact)`
- `莫娜 Mona (Genshin Impact)`
- `妮露 Nilou (Genshin Impact)`
- `荧 Player Female (Genshin Impact)`
- `埃洛伊 Aloy (Genshin Impact)`
- `柯莱 Collei (Genshin Impact)`
- `多莉 Dori (Genshin Impact)`
- `提纳里 Tighnari (Genshin Impact)`
- `砂糖 Sucrose (Genshin Impact)`
- `行秋 Xingqiu (Genshin Impact)`
- `奥兹 Oz (Genshin Impact)`
- `五郎 Gorou (Genshin Impact)`
- `达达利亚 Tartalia (Genshin Impact)`
- `七七 Qiqi (Genshin Impact)`
- `申鹤 Shenhe (Genshin Impact)`
- `莱依拉 Layla (Genshin Impact)`
- `菲谢尔 Fishl (Genshin Impact)`

### 🌸 Sanoba Witch
- `綾地 寧々 Ayachi Nene (Sanoba Witch)`
- `因幡 めぐる Inaba Meguru (Sanoba Witch)`
- `椎葉 紬 Shiiba Tsumugi (Sanoba Witch)`
- `仮屋 和奏 Kariya Wakama (Sanoba Witch)`
- `戸隠 憧子 Togakushi Touko (Sanoba Witch)`

---

## 📜 Lisensi
[MIT License](LICENSE)

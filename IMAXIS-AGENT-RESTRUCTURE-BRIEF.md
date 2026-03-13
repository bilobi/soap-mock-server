# IMAXIS Workspace — Agent & Skills Yapısal Düzenleme Brief'i

Bu brief Claude Code'a verilmek üzere hazırlanmıştır.
Adım adım uygula, her adım sonrası onay iste.

---

## Mevcut Durum

```
imaxis/                         # workspace root — CLAUDE.md YOK
├── imaxis-forge/
│   └── .claude/
│       └── skills/             # MEVCUT — taşınacak + korunacak
├── imaxis-ui/
│   └── .claude/
│       └── skills/             # MEVCUT — taşınacak + korunacak
├── imaxis-framework/           # .claude YOK
└── imaxis-go/
    └── services/
        └── *-service/          # .claude YOK
```

---

## Hedef Yapı

```
imaxis/
├── CLAUDE.md                           # workspace genel kurallar (YENİ)
├── .claude/
│   ├── skills/                         # MERKEZİ skill library (YENİ)
│   │   ├── entity-yaml/                # forge'dan taşı
│   │   ├── form-shell/                 # ui'dan taşı
│   │   ├── go-service/                 # forge'dan taşı (varsa)
│   │   └── ...diğerleri               # mevcut tüm skill'ler
│   └── agents/                         # agent tanımları (YENİ)
│       ├── architect.md
│       ├── backend-gen.md
│       ├── frontend-gen.md
│       └── migration-gen.md
│
├── imaxis-forge/
│   ├── CLAUDE.md                       # forge-spesifik context (YENİ/GÜNCELLE)
│   └── .claude/                        # KALDIR (merkeze taşındıktan sonra)
│
├── imaxis-ui/
│   ├── CLAUDE.md                       # ui-spesifik context (YENİ/GÜNCELLE)
│   └── .claude/                        # KALDIR (merkeze taşındıktan sonra)
│
├── imaxis-framework/
│   └── CLAUDE.md                       # framework context (YENİ)
│
└── imaxis-go/
    ├── CLAUDE.md                       # go services genel context (YENİ)
    └── services/
        └── *-service/
            └── CLAUDE.md               # servis-spesifik (isteğe bağlı, YENİ)
```

---

## Adım Adım Talimatlar

### ADIM 1 — Mevcut skills'leri keşfet ve listele

```
imaxis-forge/.claude/skills/ ve imaxis-ui/.claude/skills/ altındaki tüm
klasörleri ve SKILL.md dosyalarını listele. Çakışan isimler varsa raporla.
Henüz hiçbir şey taşıma, sadece raporla.
```

---

### ADIM 2 — Workspace root .claude/skills oluştur ve taşı

```
imaxis/ root'ta .claude/skills/ klasörünü oluştur.
forge ve ui'daki tüm skill klasörlerini buraya taşı.
Çakışan isim varsa içerikleri karşılaştır, hangisi daha güncel/kapsamlıysa onu al,
diğerini _deprecated prefix'iyle yanına bırak, onay iste.
```

---

### ADIM 3 — Workspace root CLAUDE.md yaz

Aşağıdaki yapıya uygun içerik oluştur:

```markdown
# IMAXIS Workspace

## Projeler
- imaxis-forge: YAML-driven codegen engine (Go). Entity şemasından Go microservice, migration, React form üretir.
- imaxis-ui: React frontend. Multi-tenant SaaS UI, shadcn/ui, dynamic i18n, OPA-based auth.
- imaxis-framework: Shared library. Multi-tenant core, OPA entegrasyonu, shared types.
- imaxis-go: Üretilmiş Go microservice'ler. Her servis services/*-service altında bağımsız.

## Genel Kurallar
- Tüm entity tanımları YAML şemasına uygun olmalı (bkz: entity-yaml skill)
- Multi-tenant: her işlemde tenant_id zorunlu
- OPA: yetkilendirme kararları OPA'ya delege edilir, servis içinde hard-code yapma
- i18n: tüm kullanıcıya dönük string'ler i18n key'i üzerinden geçmeli

## Skills
Merkezi skill library: imaxis/.claude/skills/
Görev başlamadan önce ilgili skill'i oku.

## Agent İş Akışı
Yeni entity geliştirirken sıra:
1. architect → entity YAML tasarla ve validate et
2. backend-gen → Go service üret (imaxis-go/services/)
3. frontend-gen → React form/shell üret (imaxis-ui)
4. migration-gen → DB migration yaz
```

---

### ADIM 4 — Proje bazlı CLAUDE.md'leri yaz

Her proje için ayrı CLAUDE.md oluştur. İçerik şablonu:

**imaxis-forge/CLAUDE.md**
```markdown
# imaxis-forge

Codegen engine. YAML entity şemasından kod üreten Go uygulaması.

## Kritik Dosyalar
- forge/schema/: Entity YAML şema tanımları
- forge/templates/: Go/React üretim şablonları
- forge/cmd/: CLI entry point'ler

## Bu Projede Çalışırken
- entity-yaml skill'ini her zaman yükle
- Şema değişikliği → migration-gen agent'ı tetikle
- Template değişikliği → ilgili servis re-generate edilmeli

## Bağımlılıklar
- imaxis-framework: shared types için
```

**imaxis-ui/CLAUDE.md**
```markdown
# imaxis-ui

React SaaS frontend. Multi-tenant, dynamic i18n, OPA-based permission UI.

## Kritik Dosyalar
- src/components/: shadcn/ui baz bileşenler
- src/forms/: forge'dan üretilmiş form shell'ler
- src/i18n/: dil dosyaları

## Bu Projede Çalışırken
- form-shell skill'ini yükle
- Yeni bileşen → i18n key'lerini de ekle
- Permission kontrolü → OPA policy'ye göre, hard-code etme

## Bağımlılıklar
- imaxis-framework: auth/OPA shared types
```

**imaxis-framework/CLAUDE.md**
```markdown
# imaxis-framework

Shared library. Tüm servisler bu paketi kullanır, değişiklik yaparken dikkatli ol.

## Kritik Dosyalar
- auth/: OPA client, JWT middleware
- tenant/: Multi-tenant context
- types/: Shared entity types

## Bu Projede Çalışırken
- Breaking change → tüm imaxis-go servislerini etkiler
- API değişikliği öncesi imaxis-go/services listesini kontrol et
```

**imaxis-go/CLAUDE.md**
```markdown
# imaxis-go

Üretilmiş Go microservice'ler. Çoğu imaxis-forge tarafından generate edilmiştir.
Manuel düzenleme yaparken forge şablonunu da güncelle, yoksa re-generate'de kaybolur.

## Servis Yapısı
services/*-service/:
  - main.go
  - handler/
  - repository/
  - CLAUDE.md (servis-spesifik context)

## Bu Projede Çalışırken
- go-service skill'ini yükle
- Servis port haritası: (buraya ekle)
- Her servis bağımsız binary, shared dep → imaxis-framework
```

---

### ADIM 5 — Agent tanımlarını yaz

`imaxis/.claude/agents/` altına aşağıdaki agent'ları oluştur:

**architect.md** — Entity YAML tasarımı ve validasyon
**backend-gen.md** — Go service üretimi, imaxis-go'ya yazar
**frontend-gen.md** — React form/shell üretimi, imaxis-ui'ya yazar
**migration-gen.md** — DB migration üretimi

Her agent dosyası şu yapıda olsun:
```markdown
# [Agent Adı]

## Rol
[Tek cümle görev tanımı]

## Hangi Dosyalara Bakar
[Sadece bu agent'ın erişmesi gereken path'ler]

## Hangi Skill'leri Kullanır
[İlgili skill isimleri]

## Girdi Formatı
[Bu agent'a nasıl prompt atılmalı]

## Çıktı Formatı
[Ne üretir, nereye yazar]

## Kısıtlar
[Ne yapmamalı, neye dokunmamalı]
```

---

### ADIM 6 — Eski .claude klasörlerini temizle

forge ve ui altındaki `.claude/skills/` klasörleri artık boş olmalı (taşındı).
Tamamen kaldırmadan önce onay iste.

---

## Notlar

- Her adım sonrası dur ve özet ver, devam için onay iste
- Dosya silme işlemlerinde her zaman onay iste
- Mevcut SKILL.md içeriklerini değiştirme, sadece taşı
- Çakışma veya belirsizlik varsa tahminde bulunma, sor

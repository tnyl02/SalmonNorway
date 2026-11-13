# หลังบ้านใน Docker (CON-TECH-02 "deploy ด้วย Docker", T-50 ของ 004 รับทำตาม Q-APV-38)
# build ใน image ของ Go แล้ววางโปรแกรมใน image เล็กที่ไม่มีเครื่องมืออื่น (ทีมเลือกเอง ไม่ได้มาจาก spec)
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/server /app/server
# migration และข้อมูลตั้งต้น main.go อ่านด้วย path เทียบจากโฟลเดอร์ที่รัน (db/migrations, db/seed/seed.sql)
COPY db /app/db
EXPOSE 8443
ENTRYPOINT ["/app/server"]

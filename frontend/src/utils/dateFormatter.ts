
// ฟังก์ชันสำหรับจัดรูปแบบวันที่

export function formatDateThai(dateString: string | Date): string {
  if (!dateString) return '';

  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;

  // Check if date is valid
  if (isNaN(date.getTime())) return '';

  const options: Intl.DateTimeFormatOptions = {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  };

  return date.toLocaleDateString('th-TH', options);
}

// Format date to Thai format with full month name (e.g., "28 ธันวาคม 2567")
export function formatDateThaiLong(dateString: string | Date): string {
  if (!dateString) return '';

  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;

  if (isNaN(date.getTime())) return '';

  const options: Intl.DateTimeFormatOptions = {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  };

  return date.toLocaleDateString('th-TH', options);
}

// Format date with time (e.g., "28 ธ.ค. 2567 14:30")
export function formatDateTimeThai(dateString: string | Date): string {
  if (!dateString) return '';

  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;

  if (isNaN(date.getTime())) return '';

  const options: Intl.DateTimeFormatOptions = {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  };

  return date.toLocaleDateString('th-TH', options);
}

// Format date to ISO format (YYYY-MM-DD)
export function formatDateISO(dateString: string | Date): string {
  if (!dateString) return '';

  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;

  if (isNaN(date.getTime())) return '';

  return date.toISOString().split('T')[0] ?? '';
}

// Get relative time (e.g., "2 วันที่แล้ว", "เมื่อวาน")
export function getRelativeTimeThai(dateString: string | Date): string {
  if (!dateString) return '';

  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;

  if (isNaN(date.getTime())) return '';

  const now = new Date();
  const diffTime = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return 'วันนี้';
  if (diffDays === 1) return 'เมื่อวาน';
  if (diffDays < 7) return `${diffDays} วันที่แล้ว`;
  if (diffDays < 30) return `${Math.floor(diffDays / 7)} สัปดาห์ที่แล้ว`;
  if (diffDays < 365) return `${Math.floor(diffDays / 30)} เดือนที่แล้ว`;
  return `${Math.floor(diffDays / 365)} ปีที่แล้ว`;
}

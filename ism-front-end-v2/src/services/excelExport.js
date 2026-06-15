import ExcelJS from 'exceljs';

export async function exportExcel(data, fields, fileName) {
  const workbook = new ExcelJS.Workbook();
  const worksheet = workbook.addWorksheet('数据');

  const headers = Object.keys(fields);
  
  worksheet.columns = headers.map((header, index) => ({
    header: header,
    key: header,
    width: index < 2 ? 20 : 15
  }));

  const headerRow = worksheet.getRow(1);
  headerRow.font = {
    bold: true,
    size: 11
  };
  headerRow.fill = {
    type: 'pattern',
    pattern: 'solid',
    fgColor: { argb: 'FFD3D3D3' }
  };
  headerRow.alignment = {
    horizontal: 'center',
    vertical: 'middle'
  };

  if (!data || data.length === 0) {
    return;
  }

  data.forEach((row) => {
    const rowData = {};
    headers.forEach(header => {
      const fieldConfig = fields[header];
      let value = '';
      
      if (typeof fieldConfig === 'string') {
        value = row[fieldConfig] !== undefined ? row[fieldConfig] : '';
      } else if (typeof fieldConfig === 'object' && fieldConfig.field) {
        const rawValue = row[fieldConfig.field];
        if (fieldConfig.callback && typeof fieldConfig.callback === 'function') {
          value = fieldConfig.callback(rawValue);
        } else {
          value = rawValue !== undefined ? rawValue : '';
        }
      }
      rowData[header] = value;
    });
    const newRow = worksheet.addRow(rowData);
    newRow.alignment = {
      horizontal: 'center',
      vertical: 'middle'
    };
    newRow.font = { size: 10 };
    newRow.height = 18;
  });

  worksheet.eachRow({ includeEmpty: false }, (row) => {
    row.eachCell((cell) => {
      cell.border = {
        top: { style: 'thin', color: { argb: 'FF000000' } },
        left: { style: 'thin', color: { argb: 'FF000000' } },
        bottom: { style: 'thin', color: { argb: 'FF000000' } },
        right: { style: 'thin', color: { argb: 'FF000000' } }
      };
    });
  });

  const buffer = await workbook.xlsx.writeBuffer();
  const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `${fileName}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

export async function exportExcelWithStyle(data, fields, fileName, queryDate = '', showHeader = true) {
  const workbook = new ExcelJS.Workbook();
  const worksheet = workbook.addWorksheet('报表数据');

  const headers = Object.keys(fields);
  const fieldKeys = Object.values(fields);

  const filteredHeaders = headers.filter(h => h !== '序号');
  const filteredFieldKeys = fieldKeys.filter((_, index) => headers[index] !== '序号');
  
  worksheet.columns = filteredHeaders.map((header, index) => ({
    header: header,
    key: header,
    width: index < 2 ? 15 : 10
  }));

  const headerRow = worksheet.getRow(1);
  headerRow.font = {
    bold: true,
    size: 11
  };
  headerRow.fill = {
    type: 'pattern',
    pattern: 'solid',
    fgColor: { argb: 'FFD3D3D3' }
  };
  headerRow.alignment = {
    horizontal: 'center',
    vertical: 'middle'
  };

  if (showHeader) {
    const now = new Date();
    const currentTimeStr = formatDateTime(now);
    const weekDayStr = getWeekDay(now);
    
    let reportType = '日报表';
    if (fileName.includes('差值')) reportType = '差值报表';
    else if (fileName.includes('周')) reportType = '周报表';
    else if (fileName.includes('月')) reportType = '月报表';
    else if (fileName.includes('年')) reportType = '年报表';
    else if (fileName.includes('小时')) reportType = '小时报表';

    let queryDateStr = '';
    if (queryDate) {
      if (Array.isArray(queryDate)) {
        queryDateStr = `${formatQueryDateValue(queryDate[0])} ~ ${formatQueryDateValue(queryDate[1])}`;
      } else {
        const formattedQueryDate = formatQueryDateValue(queryDate);
        const dateObj = new Date(formattedQueryDate);
        if (!isNaN(dateObj.getTime())) {
          queryDateStr = `${formattedQueryDate} (${getWeekDay(dateObj)})`;
        } else {
          queryDateStr = formattedQueryDate;
        }
      }
    }

    worksheet.spliceRows(1, 0, []);
    worksheet.spliceRows(1, 0, []);
    worksheet.spliceRows(1, 0, []);
    worksheet.spliceRows(1, 0, []);

    const infoRow = worksheet.getRow(1);
    infoRow.getCell(1).value = `查询时间: ${queryDateStr || currentTimeStr}`;
    infoRow.getCell(1).font = { size: 11 };
    infoRow.getCell(1).alignment = { horizontal: 'left', vertical: 'middle' };
    worksheet.mergeCells(1, 1, 1, Math.max(filteredHeaders.length, 1));
    
    const currentTimeRow = worksheet.getRow(2);
    const timeCell = currentTimeRow.getCell(1);
    timeCell.value = `当前时间: ${currentTimeStr} (${weekDayStr})`;
    timeCell.font = { size: 11 };
    timeCell.alignment = { horizontal: 'left', vertical: 'middle' };
    worksheet.mergeCells(2, 1, 2, Math.max(filteredHeaders.length, 1));

    const titleRow = worksheet.getRow(3);
    titleRow.getCell(1).value = reportType;
    titleRow.getCell(1).font = { bold: true, size: 16 };
    titleRow.getCell(1).alignment = { horizontal: 'left' };
    worksheet.mergeCells(3, 1, 3, Math.max(filteredHeaders.length, 1));

    if (data.length > 0 && filteredHeaders.length > 2) {
      const summaryRowData = { '设备名称': '----', '数据名称': '汇总' };
      const avgRowData = { '设备名称': '----', '数据名称': '平均值' };

      for (let i = 2; i < filteredHeaders.length; i++) {
        const colKey = filteredHeaders[i];
        let sum = 0;
        let count = 0;
        
        data.forEach(row => {
          const cellValue = row[colKey];
          const numValue = parseFloat(cellValue);
          if (!isNaN(numValue)) {
            sum += numValue;
            count++;
          }
        });

        const avg = count > 0 ? sum / count : 0;
        summaryRowData[colKey] = sum.toFixed(2);
        avgRowData[colKey] = avg.toFixed(2);
      }

      worksheet.spliceRows(7, 0, summaryRowData, avgRowData);

      const summaryRow = worksheet.getRow(7);
      summaryRow.alignment = { horizontal: 'center', vertical: 'middle' };
      summaryRow.font = { size: 10 };
      
      const avgRow = worksheet.getRow(8);
      avgRow.alignment = { horizontal: 'center', vertical: 'middle' };
      avgRow.font = { size: 10 };
    }
  }

  if (!data || data.length === 0) {
    return;
  }
  
  data.forEach((row) => {
    const rowData = [];
    filteredHeaders.forEach(header => {
      rowData.push(row[header] !== undefined ? row[header] : '');
    });
    const newRow = worksheet.addRow(rowData);
    newRow.alignment = {
      horizontal: 'center',
      vertical: 'middle'
    };
    newRow.font = { size: 10 };
    newRow.height = 18;
  });

  worksheet.eachRow({ includeEmpty: false }, (row) => {
    row.eachCell((cell) => {
      cell.border = {
        top: { style: 'thin', color: { argb: 'FF000000' } },
        left: { style: 'thin', color: { argb: 'FF000000' } },
        bottom: { style: 'thin', color: { argb: 'FF000000' } },
        right: { style: 'thin', color: { argb: 'FF000000' } }
      };
    });
  });

  const buffer = await workbook.xlsx.writeBuffer();
  const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `${fileName}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

function formatDateTime(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

function formatQueryDateValue(value) {
  if (!value) {
    return '';
  }
  if (value && typeof value.format === 'function') {
    return value.format('YYYY-MM-DD HH:mm:ss');
  }
  const date = value instanceof Date ? value : new Date(value);
  if (!isNaN(date.getTime())) {
    return formatDateTime(date);
  }
  return String(value);
}

function getWeekDay(date) {
  const weekDays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'];
  return weekDays[date.getDay()];
}

export default {
  exportExcelWithStyle
}

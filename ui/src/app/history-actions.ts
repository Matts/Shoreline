'use server';

import fs from 'fs'
import path from 'path'

export default async function getHistory() {
  const filePath = path.join(process.cwd(), '../logs', 'shoreline.log');
  const fileContent = fs.readFileSync(filePath, 'utf-8')

  console.log(filePath);

  return fileContent
}
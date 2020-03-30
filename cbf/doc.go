// package cbf implement a (MicroSoft) Compound Binary File reader.
// Note that the focus is on access to data in CBFs so stuff concerning writing CBFs is omitted.
// This package rely on the document from
// https://winprotocoldoc.blob.core.windows.net/productionwindowsarchives/SupportTech/WindowsCompoundBinaryFileFormatSpecification.pdf
// Naming will be matched as closely as possible, though hungarian notation (default MS) is omitted.
package cbf

// TODO
// 1) The master FAT is taken to be only 1o09 entries - need to handle 109+ values (Chained FATs)

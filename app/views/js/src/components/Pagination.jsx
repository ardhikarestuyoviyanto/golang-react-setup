import { CPagination, CPaginationItem } from "@coreui/react";

export default function Pagination({ currentPage, totalPage, onPageChange }) {
  const handlePrev = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPage) {
      onPageChange(currentPage + 1);
    }
  };

  return (
    <CPagination aria-label="Page navigation example" align="end">
      {/* Previous */}
      <CPaginationItem
        aria-label="Previous"
        onClick={handlePrev}
        disabled={currentPage === 1}
        style={{ cursor: currentPage === 1 ? "not-allowed" : "pointer" }}
      >
        <span aria-hidden="true">&laquo;</span>
      </CPaginationItem>

      {/* Page Numbers */}
      {Array.from({ length: totalPage }, (_, index) => index + 1).map(
        (pageNumber) => (
          <CPaginationItem
            key={pageNumber}
            active={pageNumber === currentPage}
            onClick={() => onPageChange(pageNumber)}
            style={{ cursor: "pointer" }}
          >
            {pageNumber}
          </CPaginationItem>
        )
      )}

      {/* Next */}
      <CPaginationItem
        aria-label="Next"
        onClick={handleNext}
        disabled={currentPage === totalPage}
        style={{
          cursor: currentPage === totalPage ? "not-allowed" : "pointer",
        }}
      >
        <span aria-hidden="true">&raquo;</span>
      </CPaginationItem>
    </CPagination>
  );
}

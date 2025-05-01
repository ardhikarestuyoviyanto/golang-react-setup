import {
  CAlert,
  CButton,
  CCard,
  CCardBody,
  CCardHeader,
  CCol,
  CContainer,
  CFormInput,
  CRow,
  CTable,
  CTableBody,
  CTableDataCell,
  CTableHead,
  CTableHeaderCell,
  CTableRow,
  CModal,
  CModalBody,
  CModalFooter,
  CModalHeader,
  CModalTitle,
  CForm,
  CFormTextarea,
  CFormLabel,
  CTooltip,
} from "@coreui/react";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowDown,
  faArrowUp,
  faCirclePlus,
  faFilePdf,
  faMoon,
  faPenToSquare,
  faRightFromBracket,
  faSun,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import Pagination from "../components/Pagination";
import Swal from "sweetalert2";
import { signOut, toggleDarkMode } from "../redux/slicer";
import { useFormik } from "formik";
import withReactContent from "sweetalert2-react-content";
import { useNavigate } from "react-router-dom";
const MySwal = withReactContent(Swal);

export default function Home() {
  const user = useSelector((state) => state.auth.user);
  const darkMode = useSelector((state) => state.theme.darkMode);
  const today = new Date().toISOString().slice(0, 10);
  const dispatch = useDispatch();
  const [openModal, setOpenModal] = useState(false);
  const [showPdf, setShowPdf] = useState(false);
  const [pdfUrlPrev, setUrlPdfPrev] = useState("");
  const [loading, setLoading] = useState(false);
  const [loadingForm, setLoadingForm] = useState(false);
  const navigate = useNavigate();

  const [data, setData] = useState([]);
  const [totalPage, setTotalPage] = useState(0);

  const [currentPage, setCurrentPage] = useState(1);
  const [search, setSearch] = useState("");
  const [sortField, setSortField] = useState("id");
  const [sortOrder, setSortOrder] = useState("desc");
  const [itemsPerPage] = useState(10);

  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState(today);

  const formik = useFormik({
    initialValues: {
      id: "",
      task: "",
      taskDate: "",
      attachmentFile: "",
    },
    onSubmit: async (values) => {
      setLoadingForm(true);
      const formData = new FormData();
      formData.append("id", values.id);
      formData.append("task", values.task);
      formData.append("taskDate", values.taskDate);
      formData.append("attachmentFile", values.attachmentFile);

      var method = "POST";
      var endPoint = `${import.meta.env.VITE_API_BASE_URL}/tasks`;

      if (values.id !== "") {
        method = "PUT";
        endPoint = `${import.meta.env.VITE_API_BASE_URL}/tasks/${values.id}`;
      }

      fetch(`${endPoint}`, {
        method: method,
        body: formData,
        headers: {
          Authorization: `Bearer ${user.token}`,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.success) {
            MySwal.fire({
              title: "Sukses",
              text: data.message,
              icon: "success",
              confirmButtonText: "Oke",
            }).then((result) => {
              if (result.isConfirmed) {
                formik.setValues({
                  id: "",
                  task: "",
                  taskDate: "",
                  attachmentFile: "",
                });
                fetchData();
                setOpenModal(false);
              }
            });
          } else {
            MySwal.fire({
              title: "Gagal",
              text: data.error,
              icon: "error",
              confirmButtonText: "Oke",
            });
          }
        })
        .finally(() => {
          setLoadingForm(false);
        });
    },
  });

  const handlePdf = (file) => {
    setUrlPdfPrev(`${import.meta.env.VITE_API_DOMAIN}/storage/file/${file}`);
    setShowPdf(true);
  };

  const handleEdit = async (id) => {
    fetch(`${import.meta.env.VITE_API_BASE_URL}/tasks/${id}`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${user.token}`,
        "Content-Type": "application/json",
      },
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.success) {
          formik.setValues({
            id: data.data.id,
            task: data.data.task,
            taskDate: data.data.taskDate,
            attachmentFile: "",
          });
          setOpenModal(true);
        }
      });
  };

  const handleDestroy = async (id) => {
    Swal.fire({
      title: "Hapus Data",
      text: "Yakin Pengen Dihapus ?",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya",
      cancelButtonText: "Tidak",
      reverseButtons: true,
    }).then((result) => {
      if (result.isConfirmed) {
        fetch(`${import.meta.env.VITE_API_BASE_URL}/tasks/${id}`, {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${user.token}`,
            "Content-Type": "application/json",
          },
        })
          .then((res) => res.json())
          .then((data) => {
            if (data.success) {
              MySwal.fire({
                title: "Sukses",
                text: data.message,
                icon: "success",
                confirmButtonText: "Oke",
              }).then((result) => {
                if (result.isConfirmed) {
                  fetchData();
                }
              });
            }
          });
      }
    });
  };

  const fetchData = async () => {
    try {
      setLoading(true);

      const params = new URLSearchParams({
        search,
        sortBy: sortField,
        orderBy: sortOrder,
        startDate,
        endDate,
        limit: itemsPerPage,
        offset: (currentPage - 1) * itemsPerPage,
      });

      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/tasks?${params.toString()}`,
        {
          headers: {
            Authorization: `Bearer ${user.token}`,
            "Content-Type": "application/json",
          },
        }
      );

      const result = await response.json();

      if (result.success) {
        setData(result.data.tasks);
        setTotalPage(Math.ceil(result.data.total / itemsPerPage));
      } else {
        if (response.status == 401) {
          alert("Login Expired");
          dispatch(signOut());
          window.location.href = "/";
        }
        console.error("Failed to fetch tasks", result.message);
      }
    } catch (error) {
      console.error("Error fetching tasks:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearchChange = (e) => {
    setSearch(e.target.value);
    setCurrentPage(1);
  };

  const handleStartDateChange = (e) => {
    setStartDate(e.target.value);
    setCurrentPage(1);
  };

  const handleEndDateChange = (e) => {
    setEndDate(e.target.value);
    setCurrentPage(1);
  };

  const handleSort = (field) => {
    const order = sortField === field && sortOrder === "asc" ? "desc" : "asc";
    setSortField(field);
    setSortOrder(order);
  };

  const handleLogout = () => {
    Swal.fire({
      title: "Logout Aplikasi",
      text: "Yakin Pengen Logout ?",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya",
      cancelButtonText: "Tidak",
      reverseButtons: true,
    }).then((result) => {
      if (result.isConfirmed) {
        dispatch(signOut());
        navigate("/");
      }
    });
  };

  useEffect(() => {
    fetchData();
  }, [currentPage, search, sortField, sortOrder, startDate, endDate]);

  return (
    <CContainer className="mt-5">
      <CAlert className="alert alert-success">
        Selamat Datang, {user.user.name}
      </CAlert>
      <CCard>
        <CCardHeader>
          <CButton
            style={{ float: "right" }}
            className="btn btn-primary"
            onClick={() => {
              setOpenModal(!openModal);
              formik.resetForm();
            }}
          >
            <FontAwesomeIcon icon={faCirclePlus} style={{ marginRight: 2 }} />
            Tambah
          </CButton>
          <CButton
            style={{ float: "right", marginRight: 4 }}
            className="btn btn-danger"
            onClick={() => handleLogout()}
          >
            <FontAwesomeIcon
              icon={faRightFromBracket}
              style={{ marginRight: 2 }}
            />
            Logout
          </CButton>
          <CButton
            style={{ float: "right", marginRight: 4 }}
            className="btn btn-secondary"
            onClick={() => dispatch(toggleDarkMode())}
          >
            <FontAwesomeIcon icon={darkMode ? faMoon : faSun} />
          </CButton>
        </CCardHeader>
        <CCardBody>
          <CRow className="mb-3">
            <CCol className="col-sm-4 mt-2 mb-2">
              <CFormInput
                type="date"
                value={startDate}
                onChange={handleStartDateChange}
              />
            </CCol>
            <CCol className="col-sm-4 mt-2 mb-2">
              <CFormInput
                type="date"
                value={endDate}
                onChange={handleEndDateChange}
              />
            </CCol>
            <CCol className="col-sm-4 mt-2 mb-2">
              <CFormInput
                type="text"
                placeholder="Cari Data"
                value={search}
                onChange={handleSearchChange}
              />
            </CCol>
          </CRow>

          <CTable className="table table-bordered">
            <CTableHead>
              <CTableRow>
                <CTableHeaderCell
                  scope="col"
                  style={{ width: 80, cursor: "pointer" }}
                  onClick={() => {
                    handleSort("id");
                  }}
                >
                  No
                  <FontAwesomeIcon
                    style={{ marginLeft: 4 }}
                    icon={
                      sortField === "id" && sortOrder === "asc"
                        ? faArrowUp
                        : faArrowDown
                    }
                  />
                </CTableHeaderCell>
                <CTableHeaderCell
                  scope="col"
                  style={{ cursor: "pointer" }}
                  onClick={() => {
                    handleSort("task");
                  }}
                >
                  Tugas
                  <FontAwesomeIcon
                    style={{ marginLeft: 4 }}
                    icon={
                      sortField === "task" && sortOrder === "asc"
                        ? faArrowUp
                        : faArrowDown
                    }
                  />
                </CTableHeaderCell>
                <CTableHeaderCell
                  onClick={() => {
                    handleSort("task_date");
                  }}
                  scope="col"
                  style={{ width: 150, cursor: "pointer" }}
                >
                  Tanggal
                  <FontAwesomeIcon
                    style={{ marginLeft: 4 }}
                    icon={
                      sortField === "task_date" && sortOrder === "asc"
                        ? faArrowUp
                        : faArrowDown
                    }
                  />
                </CTableHeaderCell>
                <CTableHeaderCell scope="col" style={{ width: 130 }}>
                  Aksi
                </CTableHeaderCell>
              </CTableRow>
            </CTableHead>
            <CTableBody>
              {loading ? (
                <CTableRow>
                  <CTableDataCell colSpan="4" className="text-center">
                    Loading...
                  </CTableDataCell>
                </CTableRow>
              ) : data.length > 0 ? (
                data.map((task, index) => (
                  <CTableRow key={task.id}>
                    <CTableHeaderCell scope="row">
                      {(currentPage - 1) * itemsPerPage + index + 1}
                    </CTableHeaderCell>
                    <CTableDataCell>{task.task}</CTableDataCell>
                    <CTableDataCell>{task.taskDate}</CTableDataCell>
                    <CTableDataCell>
                      <CTooltip content="Edit">
                        <CButton
                          color="primary"
                          onClick={() => {
                            handleEdit(task.id);
                          }}
                          size="sm"
                          style={{
                            marginLeft: 2,
                            marginRight: 2,
                            marginBottom: 2,
                            marginTop: 2,
                          }}
                        >
                          <FontAwesomeIcon icon={faPenToSquare} />
                        </CButton>
                      </CTooltip>

                      {task.attachmentFile != null && (
                        <CTooltip content="Pdf">
                          <CButton
                            color="success"
                            onClick={() => {
                              handlePdf(task.attachmentFile);
                            }}
                            size="sm"
                            style={{
                              marginLeft: 2,
                              marginRight: 2,
                              marginBottom: 2,
                              marginTop: 2,
                            }}
                          >
                            <FontAwesomeIcon icon={faFilePdf} />
                          </CButton>
                        </CTooltip>
                      )}
                      <CTooltip content="Hapus">
                        <CButton
                          color="danger"
                          onClick={() => {
                            handleDestroy(task.id);
                          }}
                          size="sm"
                          style={{
                            marginLeft: 2,
                            marginRight: 2,
                            marginBottom: 2,
                            marginTop: 2,
                          }}
                        >
                          <FontAwesomeIcon icon={faTrash} />
                        </CButton>
                      </CTooltip>
                    </CTableDataCell>
                  </CTableRow>
                ))
              ) : (
                <CTableRow>
                  <CTableDataCell colSpan="4" className="text-center">
                    Tidak ada data
                  </CTableDataCell>
                </CTableRow>
              )}
            </CTableBody>
          </CTable>

          <Pagination
            currentPage={currentPage}
            totalPage={totalPage}
            onPageChange={(page) => {
              setCurrentPage(page);
            }}
          />
        </CCardBody>
      </CCard>

      <CModal visible={showPdf} onClose={() => setShowPdf(false)} size="xl">
        <CModalHeader>
          <CModalTitle>Preview File</CModalTitle>
        </CModalHeader>
        <CModalBody style={{ padding: 5 }}>
          <iframe
            src={pdfUrlPrev || ""}
            width="100%"
            height="600vh"
            style={{ border: "none" }}
            title="Chrome-style PDF Preview"
          ></iframe>
        </CModalBody>
      </CModal>

      <CModal visible={openModal} onClose={() => setOpenModal(false)} size="lg">
        <CModalHeader>
          <CModalTitle>
            {formik.values.id == "" ? "Tambah Tugas" : "Edit Tugas"}
          </CModalTitle>
        </CModalHeader>
        <CForm onSubmit={formik.handleSubmit}>
          <CModalBody>
            <div className="mt-2 mb-2">
              <CFormLabel>
                Tugas <span className="text-danger">*</span>
              </CFormLabel>
              <CFormTextarea
                type="text"
                placeholder="Masukkan Tugas"
                value={formik.values.task}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                name="task"
              />
            </div>

            <div className="mt-2 mb-2">
              <CFormLabel>
                Tanggal <span className="text-danger">*</span>
              </CFormLabel>
              <CFormInput
                type="date"
                value={formik.values.taskDate}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                name="taskDate"
              />
            </div>

            <div className="mt-2 mb-2">
              <CFormLabel>Lampiran File</CFormLabel>
              <CFormInput
                type="file"
                onChange={(event) => {
                  const file = event.target.files[0];
                  formik.setFieldValue("attachmentFile", file);
                }}
                onBlur={formik.handleBlur}
                name="attachmentFile"
                text="Hanya File Pdf yang Diizinkan Untuk Diupload"
              />
            </div>
          </CModalBody>
          <CModalFooter>
            <CButton type="submit" disabled={loadingForm} color="primary">
              {formik.values.id == "" ? "Tambah" : "Edit"}
            </CButton>
          </CModalFooter>
        </CForm>
      </CModal>
    </CContainer>
  );
}

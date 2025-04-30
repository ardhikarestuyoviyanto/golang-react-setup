import {
  CCardFooter,
  CContainer,
  CForm,
  CFormInput,
  CFormLabel,
  CInputGroup,
  CInputGroupText,
} from "@coreui/react";
import { CButton, CCard, CCardBody, CCardTitle } from "@coreui/react";
import { useEffect, useState } from "react";
import { useFormik } from "formik";
import { Link, useNavigate } from "react-router-dom";
import Swal from "sweetalert2";
import withReactContent from "sweetalert2-react-content";
import { useSelector } from "react-redux";

const MySwal = withReactContent(Swal);

export default function SignUp() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const user = useSelector((state) => state.auth.user);
  const navigate = useNavigate();

  const formik = useFormik({
    initialValues: {
      name: "",
      email: "",
      password: "",
      confirmPassword: "",
      captchaId: "",
      captchaImageUrl: "",
      captchaAnswer: "",
    },
    onSubmit: (values) => {
      setLoading(true);

      const formData = new FormData();
      formData.append("email", values.email);
      formData.append("password", values.password);
      formData.append("confirmPassword", values.confirmPassword);
      formData.append("name", values.name);
      formData.append("captchaId", values.captchaId);
      formData.append("captchaAnswer", values.captchaAnswer);

      fetch(`${import.meta.env.VITE_API_BASE_URL}/signUp`, {
        method: "POST",
        body: formData,
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.success) {
            MySwal.fire({
              title: "Sukses",
              text: data.message,
              icon: "success",
              confirmButtonText: "Oke",
            }).then((res) => {
              if (res.isConfirmed) {
                formik.setValues({
                  name: "",
                  email: "",
                  password: "",
                  confirmPassword: "",
                  captchaAnswer: "",
                });
                getCaptchaId();
              }
            });
          } else {
            MySwal.fire({
              title: "Gagal",
              text: data.error || "Terjadi kesalahan",
              icon: "error",
              confirmButtonText: "Oke",
            });
          }
        })
        .finally(() => {
          getCaptchaId();
          setLoading(false);
        });
    },
  });

  const getCaptchaId = async () => {
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/captcha/generateCaptcha`,
        {
          method: "GET",
        }
      );
      const data = await res.json();
      if (data.data.captchaId) {
        formik.setFieldValue("captchaId", data.data.captchaId);
        getCaptchaImage(data.data.captchaId);
      }
    } catch (error) {
      console.error("Error generating captcha:", error);
    }
  };

  const getCaptchaImage = async (captchaId) => {
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/captcha/${captchaId}/get`,
        {
          method: "GET",
        }
      );
      const blob = await res.blob();
      const imageUrl = URL.createObjectURL(blob);
      formik.setFieldValue("captchaImageUrl", imageUrl);
    } catch (error) {
      console.error("Error fetching captcha image:", error);
    }
  };

  useEffect(() => {
    if (user != undefined) {
      navigate("/home");
    }
    getCaptchaId();
  }, []);
  return (
    <>
      <CContainer className="mt-5">
        <div className="row justify-content-center">
          <div className="col-sm-12"></div>
          <CCard style={{ width: "30rem" }}>
            <CCardBody>
              <CCardTitle>Daftar</CCardTitle>
              <hr />
              <CForm onSubmit={formik.handleSubmit}>
                <div className="mt-2 mb-2">
                  <CFormInput
                    type="text"
                    label="Nama"
                    placeholder="Masukkan Nama"
                    value={formik.values.name}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    name="name"
                  />
                </div>
                <div className="mt-2 mb-2">
                  <CFormInput
                    type="email"
                    label="Email address"
                    placeholder="Masukkan Email"
                    value={formik.values.email}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    name="email"
                  />
                </div>
                <div className="mt-2 mb-2">
                  <CFormLabel>Password</CFormLabel>
                  <CInputGroup>
                    <CFormInput
                      type={showPassword ? "text" : "password"}
                      placeholder="Masukkan Password"
                      value={formik.values.password}
                      onChange={formik.handleChange}
                      onBlur={formik.handleBlur}
                      name="password"
                    />
                    <CInputGroupText
                      onClick={() => {
                        setShowPassword(!showPassword);
                      }}
                      style={{ cursor: "pointer" }}
                    >
                      {showPassword ? "Hide" : "Show"}
                    </CInputGroupText>
                  </CInputGroup>
                </div>
                <div className="mt-2 mb-2">
                  <CFormLabel>Ulangi Password</CFormLabel>
                  <CInputGroup>
                    <CFormInput
                      type={showConfirmPassword ? "text" : "password"}
                      placeholder="Masukkan Konfirmasi Password"
                      value={formik.values.confirmPassword}
                      onChange={formik.handleChange}
                      onBlur={formik.handleBlur}
                      name="confirmPassword"
                    />
                    <CInputGroupText
                      onClick={() => {
                        setShowConfirmPassword(!showConfirmPassword);
                      }}
                      style={{ cursor: "pointer" }}
                    >
                      {showConfirmPassword ? "Hide" : "Show"}
                    </CInputGroupText>
                  </CInputGroup>
                </div>
                <div className="mt-3 mb-2">
                  {formik.values.captchaImageUrl && (
                    <>
                      <img src={formik.values.captchaImageUrl} alt="Captcha" />
                      <br />
                    </>
                  )}
                  <CFormLabel>Captcha</CFormLabel>
                  <CInputGroup>
                    <CFormInput
                      type="text"
                      name="captchaAnswer"
                      value={formik.values.captchaAnswer}
                      onChange={formik.handleChange}
                      placeholder="Masukkan captcha"
                    />
                    <CInputGroupText
                      onClick={() => {
                        getCaptchaId();
                      }}
                      style={{ cursor: "pointer" }}
                    >
                      Refresh
                    </CInputGroupText>
                  </CInputGroup>
                </div>
                <CButton
                  className="btn btn-primary mt-3"
                  type="submit"
                  disabled={loading}
                >
                  Daftar
                </CButton>
              </CForm>
            </CCardBody>
            <CCardFooter>
              <Link to="/" className="text-center">
                Login
              </Link>
            </CCardFooter>
          </CCard>
        </div>
      </CContainer>
    </>
  );
}

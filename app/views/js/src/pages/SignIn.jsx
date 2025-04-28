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
import { Link } from "react-router-dom";
import Swal from "sweetalert2";
import withReactContent from "sweetalert2-react-content";
import { useDispatch } from "react-redux";
import { signIn } from "../redux/slicer";

const MySwal = withReactContent(Swal);

export default function SignIn() {
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const dispatch = useDispatch();

  const formik = useFormik({
    initialValues: {
      captchaImageUrl: "",
      captchaId: "",
      captchaAnswer: "",
      email: "",
      password: "",
    },
    onSubmit: (values) => {
      setLoading(true);
      const formData = new FormData();
      formData.append("email", values.email);
      formData.append("password", values.password);
      formData.append("captchaId", values.captchaId);
      formData.append("captchaAnswer", values.captchaAnswer);

      fetch(`${import.meta.env.VITE_API_BASE_URL}/signIn`, {
        method: "POST",
        body: formData,
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.success) {
            dispatch(signIn(data.data));
            window.location.href = "/dashboard";
          } else {
            MySwal.fire({
              title: "Gagal",
              text: data.error || "Terjadi kesalahan",
              icon: "error",
              confirmButtonText: "Oke",
            });
          }
          getCaptchaId();
        })
        .finally(() => {
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
    getCaptchaId();
  }, []);
  return (
    <>
      <CContainer className="mt-5">
        <div className="row justify-content-center">
          <div className="col-sm-12"></div>
          <CCard style={{ width: "30rem" }}>
            <CCardBody>
              <CCardTitle> Login</CCardTitle>
              <hr />
              <CForm>
                <div className="mt-2 mb-2">
                  <CFormInput
                    type="email"
                    label="Email address"
                    placeholder="Masukkan Email"
                    name="email"
                    value={formik.values.email}
                    onChange={formik.handleChange}
                  />
                </div>
                <div className="mt-2 mb-2">
                  <CFormLabel>Password</CFormLabel>
                  <CInputGroup>
                    <CFormInput
                      type={showPassword ? "text" : "password"}
                      placeholder="Masukkan Password"
                      name="password"
                      value={formik.values.password}
                      onChange={formik.handleChange}
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
                  onClick={formik.handleSubmit}
                  className="btn btn-primary mt-3"
                  type="submit"
                  disabled={loading}
                >
                  Login
                </CButton>
              </CForm>
            </CCardBody>
            <CCardFooter>
              <Link to="/signup" className="text-center">
                Daftar
              </Link>
            </CCardFooter>
          </CCard>
        </div>
      </CContainer>
    </>
  );
}

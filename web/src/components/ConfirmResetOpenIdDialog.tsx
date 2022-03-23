import { useEffect } from "react";
import { showDialog } from "./Dialog";
import useLoading from "../hooks/useLoading";
import toastHelper from "./Toast";
import { userService } from "../services";
import "../less/confirm-reset-openid-dialog.less";

interface Props extends DialogProps {}

const ConfirmResetOpenIdDialog: React.FC<Props> = ({ destroy }: Props) => {
  const resetBtnClickLoadingState = useLoading(false);

  useEffect(() => {
    // do nth
  }, []);

  const handleCloseBtnClick = () => {
    destroy();
  };

  const handleConfirmBtnClick = async () => {
    if (resetBtnClickLoadingState.isLoading) {
      return;
    }

    resetBtnClickLoadingState.setLoading();
    try {
      await userService.resetOpenId();
    } catch (error) {
      toastHelper.error("Request reset open API failed.");
      return;
    }
    toastHelper.success("Reset open API succeeded.");
    handleCloseBtnClick();
  };

  return (
    <>
      <div className="dialog-header-container">
        <p className="title-text">Reset Open API</p>
        <button className="btn close-btn" onClick={handleCloseBtnClick}>
          <img className="icon-img" src="/icons/close.svg" />
        </button>
      </div>
      <div className="dialog-content-container">
        <p className="warn-text">
          ⚠️ The existing API will be invalidated and a new one will be generated, are you sure you want to reset?
        </p>
        <div className="btns-container">
          <span className="btn cancel-btn" onClick={handleCloseBtnClick}>
            Cancel
          </span>
          <span className={`btn confirm-btn ${resetBtnClickLoadingState.isLoading ? "loading" : ""}`} onClick={handleConfirmBtnClick}>
            Reset!
          </span>
        </div>
      </div>
    </>
  );
};

function showConfirmResetOpenIdDialog() {
  showDialog(
    {
      className: "confirm-reset-openid-dialog",
    },
    ConfirmResetOpenIdDialog
  );
}

export default showConfirmResetOpenIdDialog;

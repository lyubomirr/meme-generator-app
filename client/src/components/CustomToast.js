import { DefaultToast } from 'react-toast-notifications';
const CustomToast = ({ children, ...props }) => (
  <DefaultToast {...props}>
    <div className="custom-toast">
        {children}
    </div>
  </DefaultToast>
);

export default CustomToast

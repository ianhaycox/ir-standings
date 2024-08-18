import { useForm } from "react-hook-form";
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import { useAppDispatch } from "../../app/hooks"
import { loginAsync, Credentials } from "./loginSlice"

export const Login = () => {
    const dispatch = useAppDispatch();

    // form validation rules
    const validationSchema = Yup.object().shape({
        username: Yup.string().required('Username is required'),
        password: Yup.string().required('Password is required')
    });
    const formOptions = { resolver: yupResolver(validationSchema) };

    // get functions to build form with useForm() hook
    const { register, handleSubmit, formState } = useForm(formOptions);
    const { errors, isSubmitting } = formState;

    function onSubmit(data:Credentials) {
        return dispatch(loginAsync(data));
    }

    return (
        <div className="card m-3">
            <h4 className="card-header">Login</h4>
            <div className="card-body">
                <form onSubmit={handleSubmit(onSubmit)}>
                    <div className="mb-3">
                        <label htmlFor="username" className="form-label">Username</label>
                        <input id="username" type="text" autoComplete="username" {...register('username')} className={`form-control ${errors.username ? 'is-invalid' : ''}`} />
                        <div className="invalid-feedback">{errors.username?.message}</div>
                    </div>
                    <div className="mb-3">
                        <label  htmlFor="password" className="form-label">Password</label>
                        <input id="password" type="password" autoComplete="current-password" {...register('password')} className={`form-control ${errors.password ? 'is-invalid' : ''}`} />
                        <div className="invalid-feedback">{errors.password?.message}</div>
                    </div>
                    <button disabled={isSubmitting} className="btn btn-primary">
                        {isSubmitting && <span className="spinner-border spinner-border-sm me-1"></span>}
                        Login
                    </button>
                </form>
            </div>
        </div>
    )
}

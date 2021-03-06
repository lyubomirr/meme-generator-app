import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { getCurrentUser } from '../auth';

const PrivateRoute = ({component: Component, ...rest}) => {
    return (
        <Route {...rest} render={props => (
            getCurrentUser() ?
                <Component {...props} />
            : <Redirect to="/login" />
        )} />
    );
};

export default PrivateRoute;
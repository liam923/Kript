import Aux from '../../../hoc/Aux/Aux';
import React from 'react';

export class Password {
  static card(props) {
    return <Aux>
      <h4>{props.datum.url}</h4>
      <p>{props.datum.username}</p>
    </Aux>;
  }

  static navItem(props) {

  }
}

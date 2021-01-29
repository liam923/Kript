import React from 'react';

import Button from '../common/Button/Button';
import NavigationItem from './NavigationItem/NavigationItem';

import styles from './NavigationBar.module.css';

const navigationBar = (props) => (
    <div className={styles.NavigationBar}>
        <div>
            <NavigationItem navType='all' clicked={() => props.itemClicked('')} active={props.filter === ''}/>
            <NavigationItem navType='passwords' clicked={() => props.itemClicked('passwords')} active={props.filter === 'passwords'}/>
            <NavigationItem navType='creditcards' clicked={() => props.itemClicked('creditcards')} active={props.filter === 'creditcards'}/>
            <NavigationItem navType='notes' clicked={() => props.itemClicked('notes')} active={props.filter === 'notes'}/>
            <NavigationItem navType='codes' clicked={() => props.itemClicked('codes')} active={props.filter === 'codes'}/>
        </div>
        <Button clicked={props.addDataClicked}>Add Data</Button>
        <Button clicked={props.logOutClicked}>Log Out</Button>
    </div>
);

export default navigationBar;
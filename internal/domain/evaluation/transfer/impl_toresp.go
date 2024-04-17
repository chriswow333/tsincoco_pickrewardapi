package transfer

// payDTO "pickrewardapi/internal/app/pay/dto"

// func (im *impl) transferToPayloadRespDTO(ctx context.Context, payload *domain.Payload, evaluationRespDTO *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferToPayloadRespDTO]"

// 	if payload == nil {
// 		log.WithFields(log.Fields{
// 			"pos":        logPos,
// 			"payload.ID": payload.ID,
// 		}).Error("payload is nil")
// 		return errors.New("payload is nil")
// 	}

// 	if payload.PayloadType == domain.SelfPayloadType {
// 		if payload.Payloads == nil {
// 			log.WithFields(log.Fields{
// 				"pos":        logPos,
// 				"payload.ID": payload.ID,
// 			}).Error("payloads is nil")
// 			return errors.New("payloads is nil")
// 		}

// 		for _, p := range payload.Payloads {
// 			if err := im.transferToPayloadRespDTO(ctx, p, evaluationRespDTO); err != nil {
// 				log.WithFields(log.Fields{
// 					"pos":        logPos,
// 					"payload.ID": payload.ID,
// 				}).Error("transferToPayloadRespDTO failed ", err)
// 				return err
// 			}

// 		}
// 		return nil
// 	}

// 	if payload.PayloadType == domain.ContainerPayloadType {
// 		if err := im.transferToContainerRespDTO(ctx, payload.Container, evaluationRespDTO); err != nil {
// 			log.WithFields(log.Fields{
// 				"pos":        logPos,
// 				"payload.ID": payload.ID,
// 			}).Error("transferToContainerRespDTO failed ", err)
// 			return err
// 		}
// 	}

// 	return nil

// }

// func (im *impl) transferToContainerRespDTO(ctx context.Context, container *domain.Container, evaluationRespDTO *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferToContainerRespDTO]"

// 	if container == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("container is nil")
// 		return errors.New("container is nil")
// 	}

// 	if container.ContainerType == domain.InnerContainer {
// 		if container.InnerContainers == nil {
// 			log.WithFields(log.Fields{
// 				"pos": logPos,
// 			}).Error("innerContainers is nil")
// 			return errors.New("innerContainers is nil")
// 		}
// 		for _, c := range container.InnerContainers {
// 			im.transferToContainerRespDTO(ctx, c, evaluationRespDTO)
// 		}
// 		return nil
// 	}

// 	if container.ContainerType == domain.ChannelContainer {
// 		if err := im.transferChannelEvaluationRespDTO(ctx, container, evaluationRespDTO); err != nil {
// 			log.WithFields(log.Fields{
// 				"pos": logPos,
// 			}).Error("transferChannelEvaluationRespDTO failed")
// 			return err
// 		}
// 		return nil
// 	}

// 	if container.ContainerType == domain.PayContainer {
// 		if err := im.transferPayEvaluationRespDTO(ctx, container, evaluationRespDTO); err != nil {
// 			log.WithFields(log.Fields{
// 				"pos": logPos,
// 			}).Error("transferPayEvaluationRespDTO failed")
// 			return err
// 		}
// 		return nil
// 	}

// 	if container.ContainerType == domain.CardRewardTaskLabelContainer {
// 		if err := im.transferTaskEvaluationRespDTO(ctx, container, evaluationRespDTO); err != nil {
// 			log.WithFields(log.Fields{
// 				"pos": logPos,
// 			}).Error("transferTaskEvaluationRespDTO failed")
// 			return err
// 		}
// 		return nil
// 	}

// 	if container.ContainerType == domain.ConstraintContainer {
// 		if err := im.transferConstraintsEvaluationRespDTO(ctx, container, evaluationRespDTO); err != nil {
// 			log.WithFields(log.Fields{
// 				"pos": logPos,
// 			}).Error("transferConstraintsEvaluationRespDTO failed")
// 			return err
// 		}
// 		return nil
// 	}

// 	if container.ContainerType == domain.ChannelLabelContainer {
// 		if err := im.transferChannelLabelEvaluationRespDTO(ctx, container, evaluationRespDTO); err != nil {
// 			log.WithFields(log.Fields{
// 				"pos": logPos,
// 			}).Error("transferLabelEvaluationRespDTO failed")
// 			return err
// 		}
// 		return nil
// 	}

// 	return nil
// }

// func (im *impl) transferChannelEvaluationRespDTO(ctx context.Context, container *domain.Container, evaluationRespDTO *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferChannelEvaluationRespDTO]"

// 	if container.ChannelEvaluations == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("channelDTOs are nil")
// 		return errors.New("channelDTOs are nil")
// 	}

// matches := map[int32]map[string]*.ChannelDTO{}

// if container.ContainerOperator == domain.AndContainer || container.ContainerOperator == domain.OrContainer {
// 	for _, c := range container.ChannelDTOs {
// 		if _, ok := matches[c.ChannelCategoryType]; !ok {
// 			matches[c.ChannelCategoryType] = make(map[string]*channelDTO.ChannelDTO)
// 		}
// 		matches[c.ChannelCategoryType][c.ID] = c
// 	}
// }

// 	misMatches := map[int32]map[string]*channelDTO.ChannelDTO{}

// 	if container.ContainerOperator == domain.NotContainer {
// 		for _, c := range container.ChannelDTOs {
// 			if _, ok := misMatches[c.ChannelCategoryType]; !ok {
// 				misMatches[c.ChannelCategoryType] = make(map[string]*channelDTO.ChannelDTO)
// 			}
// 			misMatches[c.ChannelCategoryType][c.ID] = c
// 		}
// 	}

// 	if evaluationRespDTO.ChannelsEvaluationResp == nil {
// 		evaluationRespDTO.ChannelsEvaluationResp = &evaluationDTO.ChannelsEvaluationResp{}
// 	}

// 	if evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper == nil {
// 		evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper = make(map[int32]*evaluationDTO.ChannelEvaluationRespDTO)
// 	}

// 	for category, channels := range matches {
// 		if _, ok := evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper[category]; !ok {
// 			evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper[category] = &evaluationDTO.ChannelEvaluationRespDTO{
// 				ChannelCategoryType: category,
// 				Matches:             make(map[string]*channelDTO.ChannelDTO),
// 				MisMatches:          make(map[string]*channelDTO.ChannelDTO),
// 			}
// 		}

// 		for id, c := range channels {
// 			evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper[category].Matches[id] = c
// 		}

// 	}

// 	for category, channels := range misMatches {
// 		if _, ok := evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper[category]; !ok {
// 			evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper[category] = &evaluationDTO.ChannelEvaluationRespDTO{
// 				ChannelCategoryType: category,
// 				Matches:             make(map[string]*channelDTO.ChannelDTO),
// 				MisMatches:          make(map[string]*channelDTO.ChannelDTO),
// 			}
// 		}

// 		for id, c := range channels {
// 			evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper[category].MisMatches[id] = c
// 		}

// 	}

// 	return nil

// }

// func (im *impl) transferPayEvaluationRespDTO(ctx context.Context, container *domain.Container, evaluationRespDTO *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferPayEvaluationRespDTO]"

// 	if container.PayDTOs == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("payDTOs is nil")
// 		return errors.New("payDTOs is nil")
// 	}

// 	matches := map[string]*payDTO.PayDTO{}
// 	if container.ContainerOperator == domain.AndContainer || container.ContainerOperator == domain.OrContainer {
// 		for _, p := range container.PayDTOs {
// 			matches[p.ID] = p
// 		}
// 	}

// 	misMatches := map[string]*payDTO.PayDTO{}
// 	if container.ContainerOperator == domain.NotContainer {
// 		for _, p := range container.PayDTOs {
// 			misMatches[p.ID] = p
// 		}
// 	}

// 	if evaluationRespDTO.PayEvaluationResp == nil {
// 		evaluationRespDTO.PayEvaluationResp = &evaluationDTO.PayEvaluationRespDTO{
// 			Matches:    make(map[string]*payDTO.PayDTO),
// 			MisMatches: make(map[string]*payDTO.PayDTO),
// 		}
// 	}

// 	for id, p := range matches {
// 		evaluationRespDTO.PayEvaluationResp.Matches[id] = p
// 	}

// 	for id, p := range misMatches {
// 		evaluationRespDTO.PayEvaluationResp.MisMatches[id] = p
// 	}

// 	return nil

// }

// func (im *impl) transferCardRewardTaskLabelEvaluationRespDTO(ctx context.Context, container *domain.Container, evaluationResp *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferCardRewardTaskLabelEvaluationRespDTO]"

// 	if container.CardRewardTaskLabels == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("CardRewardTaskLabels are nil")
// 		return errors.New("CardRewardTaskLabels are nil")
// 	}

// 	return nil
// }

// func (im *impl) transferConstraintsEvaluationRespDTO(ctx context.Context, container *domain.Container, evaluationResp *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferConstraintsEvaluationRespDTO]"

// 	if container.Constraints == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("constraints is nil")
// 		return errors.New("constraints is nil")
// 	}

// 	if evaluationResp.ConstraintsEvaluationResp == nil {
// 		evaluationResp.ConstraintsEvaluationResp = &evaluationDTO.ConstraintsEvaluationRespDTO{
// 			Matches:    make(map[int32]*evaluationDTO.ConstraintDTO),
// 			MisMatches: make(map[int32]*evaluationDTO.ConstraintDTO),
// 		}
// 	}

// 	if container.ContainerOperator == domain.AndContainer || container.ContainerOperator == domain.OrContainer {
// 		for _, c := range container.Constraints {
// 			evaluationResp.ConstraintsEvaluationResp.Matches[int32(c.ConstraintType)] = &evaluationDTO.ConstraintDTO{
// 				ConstraintType: int32(c.ConstraintType),
// 				ConstraintName: c.ConstraintName,
// 				WeekDays:       c.WeekDays,
// 			}
// 		}
// 		return nil
// 	}

// 	if container.ContainerOperator == domain.NotContainer {
// 		for _, c := range container.Constraints {
// 			evaluationResp.ConstraintsEvaluationResp.MisMatches[int32(c.ConstraintType)] = &evaluationDTO.ConstraintDTO{
// 				ConstraintType: int32(c.ConstraintType),
// 				ConstraintName: c.ConstraintName,
// 				WeekDays:       c.WeekDays,
// 			}
// 		}
// 		return nil
// 	}

// 	return nil
// }

// func (im *impl) transferChannelLabelEvaluationRespDTO(ctx context.Context, container *domain.Container, evaluationResp *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferChannelLabelEvaluationRespDTO]"

// 	if container.ChannelLabelDTOs == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("ChannelLabels are nil")
// 		return errors.New("ChannelLabels are nil")
// 	}

// 	if evaluationResp.ChannelLabelEvaluationResp == nil {
// 		evaluationResp.ChannelLabelEvaluationResp = &evaluationDTO.ChannelLabelEvaluationResp{
// 			Matches:    make(map[int32]*channelDTO.ChannelLabelDTO),
// 			MisMatches: make(map[int32]*channelDTO.ChannelLabelDTO),
// 		}
// 	}

// 	if container.ContainerOperator == domain.AndContainer || container.ContainerOperator == domain.OrContainer {
// 		for _, l := range container.ChannelLabelDTOs {
// 			evaluationResp.ChannelLabelEvaluationResp.Matches[int32(l.Label)] = l
// 		}
// 		return nil
// 	}

// 	if container.ContainerOperator == domain.NotContainer {
// 		for _, l := range container.ChannelLabelDTOs {
// 			evaluationResp.ChannelLabelEvaluationResp.MisMatches[int32(l.Label)] = l
// 		}
// 		return nil
// 	}

// 	return nil
// }

// func (im *impl) transferChannelCategoryTypes(ctx context.Context, evaluationRespDTO *evaluationDTO.EvaluationRespDTO) error {
// 	logPos := "[evaluation.app.transfer][transferChannelCategoryTypes]"

// 	evaluationRespDTO.ChannelCategoryTypes = []*channelDTO.ChannelCategoryTypeDTO{}

// 	if evaluationRespDTO.ChannelsEvaluationResp != nil && evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper != nil {
// 		for _, channelRespDTO := range evaluationRespDTO.ChannelsEvaluationResp.ChannelEvaluationRespMapper {

// 			// TODO only get matches
// 			if len(channelRespDTO.Matches) != 0 {
// 				channelCategoryType, err := im.channelAppService.GetChannelCategoryTypeByCategoryType(ctx, channelRespDTO.ChannelCategoryType)
// 				if err != nil {
// 					log.WithFields(log.Fields{
// 						"pos": logPos,
// 					}).Error("channelAppService.GetChannelCategoryTypeByCategoryType is nil")
// 					return err
// 				}
// 				evaluationRespDTO.ChannelCategoryTypes = append(evaluationRespDTO.ChannelCategoryTypes, channelCategoryType)
// 			}
// 		}
// 	}

// 	return nil
// }
